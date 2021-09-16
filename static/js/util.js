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
    const payload = {'iat': timeStamp, 'exp': timeStamp + exp}
    const payloadJson = JSON.stringify(payload)
    const secret = getSecret()
    const secretHex = enSha256Hex(secret)
    const jwt = KJUR.jws.JWS.sign("HS256", headerJson, payloadJson, {hex: secretHex})
    return jwt
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