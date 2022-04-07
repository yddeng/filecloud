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