import request from '@/utils/request'

export function mkdir (parameter) {
  return request({
    url: "path/mkdir",
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