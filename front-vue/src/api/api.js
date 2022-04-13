import request from '@/utils/request'


/**
 * login func
 * parameter: {
 *     username: '',
 *     password: '',
 *     remember_me: true,
 *     captcha: '12345'
 * }
 * @param parameter
 * @returns {*}
 */
export function mkdir (parameter) {
  return request({
    url: "file/mkdir",
    method: 'post',
    data: parameter
  })
}

export function list (parameter) {
  return request({
    url: "file/list",
    method: 'post',
    data: parameter
  })
}

export function fileRemove (parameter) {
  return request({
    url: "file/remove",
    method: 'post',
    data: parameter
  })
}

export function rename (parameter) {
  return request({
    url: "file/rename",
    method: 'post',
    data: parameter
  })
}

export function mvcp (parameter) {
  return request({
    url: "file/mvcp",
    method: 'post',
    data: parameter
  })
}

export function uploadCheck (parameter) {
  return request({
    url: "upload/check",
    method: 'post',
    data: parameter
  })
}

export function uploadFile (parameter) {
  return request({
    url: "upload/upload",
    method: 'post',
    data: parameter
  })
}

export function downloadFile (parameter) {
  return request({
    url: "file/download",
    method: 'post',
    responseType:"blob",
    data: parameter
  })
}

export function sharedCreate (parameter) {
  return request({
    url: "shared/create",
    method: 'post',
    data: parameter
  })
}

export function sharedList (parameter) {
  return request({
    url: "shared/create",
    method: 'post',
    data: parameter
  })
}