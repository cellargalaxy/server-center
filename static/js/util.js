function enSha256Hex(text) {
    if (text === undefined || text == null) {
        text = ''
    }
    const hash = sha256.create()
    hash.update(text)
    return hash.hex()
}

const secretKey = 'secret'

function getSecret() {
    return localStorage.getItem(secretKey)
}

function setSecret(secret) {
    localStorage.setItem(secretKey, secret)
}

const tokenExpKey = 'tokenExp'

function getTokenExp() {
    const exp = localStorage.getItem(tokenExpKey)
    if (isNum(exp)) {
        return parseInt(exp)
    }
    return 3
}

function setTokenExp(exp) {
    localStorage.setItem(tokenExpKey, exp)
}

function enJwt() {
    const timeStamp = getTimeStamp()
    const exp = getTokenExp()
    const header = {'typ': 'JWT', 'alg': 'HS256'}
    const headerJson = JSON.stringify(header)
    const payload = {'iat': timeStamp, 'exp': timeStamp + exp, 'allow_re_request': false, 'request_id': genId()}
    const payloadJson = JSON.stringify(payload)
    const secret = getSecret()
    const secretHex = enSha256Hex(secret)
    return KJUR.jws.JWS.sign("HS256", headerJson, payloadJson, {hex: secretHex})
}

function isNum(s) {
    if (s != null && s !== '') {
        return !isNaN(s)
    }
    return false
}

function getTimeStamp() {
    return Date.parse(new Date()) / 1000
}

function formatTimestamp(timestamp, fmt) {
    return formatDate(new Date(timestamp), fmt)
}

function formatDate(date, fmt) {
    let o = {
        'M+': date.getMonth() + 1, //月份
        'D+': date.getDate(), //日
        'H+': date.getHours(), //小时
        'm+': date.getMinutes(), //分
        's+': date.getSeconds(), //秒
        'q+': Math.floor((date.getMonth() + 3) / 3), //季度
        'S': date.getMilliseconds() //毫秒
    }
    if (/(Y+)/.test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
    }
    for (let k in o) {
        if (new RegExp('(' + k + ')').test(fmt)) {
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (('00' + o[k]).substr(('' + o[k]).length)))
        }
    }
    return fmt
}

function genId() {
    return formatDate(new Date(), 'YYMMDDHHmmssS') + 'web' + randomWord()
}

function randomWord() {
    let words = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
        'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
        'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z']
    return words[Math.floor(Math.random() * words.length)]
}

function getQueryString(name) {
    let reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i")
    let r = window.location.search.substr(1).match(reg)
    if (r != null) {
        return decodeURIComponent(r[2])
    }
    return null
}

function startsWith(string, start) {
    return string.indexOf(start) === 0
}

function containWith(string, sub) {
    return string.indexOf(sub) >= 0
}

function isNum(s) {
    if (s != null && s !== '') {
        return !isNaN(s)
    }
    return false
}

function toString(value) {
    if (value === null || typeof value === 'undefined') {
        return ''
    } else if (value instanceof Object) {
        return Object.keys(value)
            .sort()
            .map(key => toString(value[key]))
            .join(' ')
    } else {
        return String(value)
    }
}

function getAddYearDate(year) {
    let before = new Date()
    before.setFullYear(before.getFullYear() + year)
    return before
}

function parse2BeijingTimestamp(date) {
    let timestamp = 0
    if (date !== undefined && date != null && date !== '') {
        timestamp = Date.parse(date + ' GMT+8') / 1000
    }
    if (!isNum(timestamp)) {
        timestamp = 0
    }
    return timestamp
}


function log(...data) {
    console.log('log', data)
}

function sortCompare(aRow, bRow, key, sortDesc, formatter, compareOptions, compareLocale) {
    let a = aRow[key]
    let b = bRow[key]
    if (isNum(a) && isNum(b)) {
        a = parseFloat(a)
        b = parseFloat(b)
        return a < b ? -1 : a > b ? 1 : 0
    }
    return toString(a).localeCompare(toString(b), compareLocale, compareOptions)
}