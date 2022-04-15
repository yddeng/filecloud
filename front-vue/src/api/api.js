import request from '@/utils/request'


export function login (parameter) {
  return request({
    url: "auth/login",
    method: 'post',
    data: parameter
  })
}


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

export function sharedList (parameter) {
  return request({
    url: "shared/list",
    method: 'post',
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

export function sharedInfo (parameter) {
  return request({
    url: "shared/info",
    method: 'post',
    data: parameter
  })
}

export function sharedPath (parameter) {
  return request({
    url: "shared/path",
    method: 'post',
    data: parameter
  })
}

export function sharedDownload (parameter) {
  return request({
    url: "shared/download",
    method: 'post',
    responseType:"blob",
    data: parameter
  })
}

export function sharedCancel (parameter) {
  return request({
    url: "shared/cancel",
    method: 'post',
    data: parameter
  })
}