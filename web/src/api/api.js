import requset from '@/utils/request';

export function getStats(query) {
    return requset({
        url: '/stats',
        method:'get',
        params: query,
    })
}

export function getPlayers(query) {
    return requset({
        url: '/players',
        method:'get',
        params: query,
    })
}