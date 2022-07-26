package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"
)

type Ntof struct {
	token  string
	cron   *cron.Cron
	client *Client
	db     *mongo.Client
}

func newNtof() *Ntof {
	db, err := newDB()
	if err != nil {
		log.Fatalln(err)
	}

	return &Ntof{
		cron:   cron.New(cron.WithSeconds()),
		client: NewClient(),
		db:     db,
	}
}

const (
	GoodSIdShangWu = 2
	GoodSIdXiaWu   = 3

	BaseURL     = "http://h5.lingyuan888.com/index.php/WebApi"
	LoginURL    = "/User/login"
	GoodListURL = "/Goods/goods"
	BuyURL      = "/Order/buy"
)

func (n *Ntof) Login(username, password string) error {
	params := map[string]string{
		"uphone":   username,
		"password": password,
	}

	out := struct {
		Token    string
		NickName string
	}{}
	err := n.client.PostForm(BaseURL+LoginURL, params, &out)
	if err != nil {
		return err
	}

	fmt.Println("login success, token:", out.Token)

	n.token = out.Token
	return nil
}

type Good struct {
	Id            string `json:"id"`
	GNum          string `json:"gnum"`
	GName         string `json:"gname"`
	OriginalPrice string `json:"orignialprice"`
	CurPrice      string `json:"curprice"`
	CTime         string `json:"ctime"`
	CId           string `json:"cid"`
	GDes          string `json:"gdes"`
	GStatus       string `json:"gstatus"` // 1 卖光了
	OwnUname      string `json:"ownuname"`
	OwnUID        string `json:"ownuid"`
	OnSale        string `json:"onsale"`
	SName         string `json:"sname"`
	Status        string `json:"status"`
}

func (n *Ntof) GoodList(page, id int) (int, []*Good, error) {
	var data struct {
		Total  string  `json:"total"`
		Offset int     `json:"offset"`
		Goods  []*Good `json:"goods"`
	}
	queryUrl := fmt.Sprintf(BaseURL+GoodListURL+"?page=%d&count=100&sid=%d&token=%s", page, id, n.token)
	err := n.client.Get(queryUrl, &data)
	if err != nil {
		return 0, nil, err
	}
	total, _ := strconv.ParseInt(data.Total, 10, 64)
	return int(total), data.Goods, nil
}

func (n *Ntof) Buy(showid, gid string, token string) error {
	params := map[string]string{
		"gid":    gid,
		"showid": showid,
		"token":  token,
	}
	var out interface{}
	err := n.client.PostForm(BaseURL+BuyURL, params, out)
	if err != nil {
		return err
	}

	fmt.Println("buy response", out)
	return nil
}

var ntof *Ntof

type player struct {
	Id        string
	Name      string
	Count     int
	Asset     float64
	Date      string
	CreatedAt time.Time
}

func main() {
	ntof = newNtof()
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	err := ntof.Login(username, password)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		allGoods   []*Good
		initialCap float64
		marketCap  float64
		players    = make(map[string]*player)
		filtter    = make(map[string]struct{})
	)

	for _, show := range []int{GoodSIdShangWu, GoodSIdXiaWu} {
		var page int
		var getGoods []*Good

		for {
			page++
			total, goods, err := ntof.GoodList(page, show)
			if err != nil {
				log.Println(err)
				return
			}

			getGoods = append(getGoods, goods...)
			allGoods = append(allGoods, goods...)

			if len(getGoods) >= total {
				break
			}
		}
	}

	for _, good := range allGoods {
		_, fok := filtter[good.Id]
		if fok {
			continue
		}
		filtter[good.Id] = struct{}{}

		_, ok := players[good.OwnUname]
		if !ok {
			players[good.OwnUname] = &player{
				Id:   good.OwnUID,
				Name: good.OwnUname,
				Date: time.Now().Format("20060102"),
			}
		}

		if good.OwnUname == "kin01" {
			fmt.Println(good)
		}

		ori, _ := strconv.ParseFloat(good.OriginalPrice, 10)
		cur, _ := strconv.ParseFloat(good.CurPrice, 10)
		marketCap += cur
		initialCap += ori
		players[good.OwnUname].Asset += cur
		players[good.OwnUname].Count++
	}

	var finalPLayers []*player
	for _, player := range players {
		finalPLayers = append(finalPLayers, player)
		//fmt.Println(fmt.Sprintf("%s: %d:  %.2f", name, player.Count, player.Asset))
	}

	UpsertPlayer(context.Background(), Players(ntof.db), finalPLayers)

	sort.Slice(finalPLayers, func(i, j int) bool {
		return finalPLayers[i].Asset > finalPLayers[j].Asset
	})

	for _, player := range finalPLayers {
		fmt.Println(fmt.Sprintf("%s: %d:  %.2f", player.Name, player.Count, player.Asset))
	}

	fmt.Println(fmt.Sprintf("Initial: %.2f", initialCap))
	fmt.Println(fmt.Sprintf("MarketCap: %.2f", marketCap))
	fmt.Println("Players: ", len(players))
	fmt.Println("Goods Count: ", len(allGoods))

	UpsertStats(context.Background(), Stats(ntof.db), &stats{
		InitialCap: initialCap,
		MarketCap:  marketCap,
		Players:    len(players),
		GoodCount:  len(allGoods),
		Date:       time.Now().Format("20060102"),
		CreatedAt:  time.Now(),
	})

	if os.Getenv("RUN") == "1" {
		runJob()
	}

	//ntof.cron.AddFunc("55 29 11 * * *", runJob)
	//ntof.cron.AddFunc("55 59 14 * * *", runJob)
	//ntof.cron.Start()
	//defer ntof.cron.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	select {
	case sig := <-sigChan:
		fmt.Println("received shutdown", "signal", sig)
	}

	fmt.Println("Graceful shutdown successful")
}

func runJob() {
	username := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	err := ntof.Login(username, password)
	if err != nil {
		log.Fatalln(err)
	}

	var start = time.Now()

	sid := GoodSIdShangWu
	hour := time.Now().Hour()
	if hour > 12 {
		sid = GoodSIdXiaWu
	}

	_, goods, err := ntof.GoodList(1, sid)
	if err != nil {
		log.Println(err)
		return
	}

	sort.Slice(goods, func(i, j int) bool {
		cur1, _ := strconv.ParseFloat(goods[i].CurPrice, 10)
		cur2, _ := strconv.ParseFloat(goods[j].CurPrice, 10)
		return cur1 > cur2
	})

	for {
		for _, good := range goods {
			// 卖光了
			if good.GStatus != "7" {
				continue
			}

			go func() {
				err := ntof.Buy(strconv.Itoa(sid), good.Id, ntof.token)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}

		time.Sleep(50 * time.Millisecond)
		if time.Now().After(start.Add(1 * time.Minute)) {
			fmt.Println("timeout exit")
			break
		}
	}
}
