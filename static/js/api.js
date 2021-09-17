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
    if (code === undefined || code == null || code === '') {
        dealErr('code为空')
        return null
    }
    start_time = parse2BeijingTimestamp(start_time)
    if (!isNum(start_time)) {
        start_time = 0
    }
    end_time = parse2BeijingTimestamp(end_time)
    if (!isNum(end_time)) {
        end_time = 0
    }

    let url = '../../api/analyzeFund'
    if (document.domain === 'localhost') {
        url += '.json'
    }
    try {
        let response = await instance.get(url, {
            params: {code: code, start_time: start_time, end_time: end_time,},
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