const instance = axios.create({timeout: 60 * 1000})
instance.interceptors.request.use(
    config => {
        config.headers['Authorization'] = 'Bearer ' + enJwt()
        return config
    },
    error => Promise.reject(error))

async function ping() {
    let url = '../../ping'
    if (document.domain === 'localhost') {
        url += '.json'
    }
    try {
        let response = await instance.get(url, {
            params: {},
            paramsSerializer: params => {
                return Qs.stringify(params, {indices: false})
            }
        })
        return dealResponse(response)
    } catch (error) {
        dealErr(error)
    }
    return null
}

async function addServerConf(server_name, version, remark, conf_text) {
    if (server_name === undefined || server_name == null || server_name === '') {
        dealErr('server_name为空')
        return null
    }
    if (version === undefined || version == null || version === '') {
        dealErr('version为空')
        return null
    }
    if (remark === undefined || remark == null || remark === '') {
        dealErr('remark为空')
        return null
    }

    let url = '../../api/addServerConf'
    if (document.domain === 'localhost') {
        url += '.json'
    }
    try {
        let response = await instance.post(url, {
            server_name: server_name,
            version: version,
            remark: remark,
            conf_text: conf_text,
        })
        return dealResponse(response)
    } catch (error) {
        dealErr(error)
    }
    return null
}

async function getLastServerConf(server_name, version) {
    if (server_name === undefined || server_name == null || server_name === '') {
        dealErr('server_name为空')
        return null
    }
    if (version === undefined || version == null || version === '') {
        dealErr('version为空')
        return null
    }

    let url = '../../api/getLastServerConf'
    if (document.domain === 'localhost') {
        url += '.json'
    }
    try {
        let response = await instance.get(url, {
            params: {server_name: server_name, version: version},
            paramsSerializer: params => {
                return Qs.stringify(params, {indices: false})
            }
        })
        return dealResponse(response)
    } catch (error) {
        dealErr(error)
    }
    return null
}

async function getLastServerConfVersion(server_name) {
    let conf = getLastServerConf(server_name, 0)
    if (conf === undefined || conf == null) {
        return 0
    }
    if (!isNaN(conf.version)) {
        return 0
    }
    return conf.version
}

async function listServerConf(server_name, version) {
    if (server_name === undefined || server_name == null || server_name === '') {
        dealErr('server_name为空')
        return null
    }
    if (version === undefined || version == null || version === '') {
        dealErr('version为空')
        return null
    }

    let url = '../../api/listServerConf'
    if (document.domain === 'localhost') {
        url += '.json'
    }
    try {
        let response = await instance.get(url, {
            params: {server_name: server_name, version: version},
            paramsSerializer: params => {
                return Qs.stringify(params, {indices: false})
            }
        })
        return dealResponse(response)
    } catch (error) {
        dealErr(error)
    }
    return null
}

async function listAllServerName() {
    let url = '../../api/listAllServerName'
    if (document.domain === 'localhost') {
        url += '.json'
    }
    try {
        let response = await instance.get(url, {
            params: {},
            paramsSerializer: params => {
                return Qs.stringify(params, {indices: false})
            }
        })
        return dealResponse(response)
    } catch (error) {
        dealErr(error)
    }
    return null
}

function dealResponse(response) {
    let result = response.data
    if (result.code !== 1) {
        dealErr(result.message)
        return null
    }
    return result.data
}

function dealErr(error) {
    let msg = JSON.stringify(error)
    if (msg === undefined || msg == null || msg === '' || msg === '{}' || msg === '[]') {
        msg = error
    }
    alert("error: " + msg)
    log(msg)
}