import { useGetData, usePostData, usePost, useGet } from './hook'

// ============ 用户管理 API ============

// 获取用户列表
export function userList(params) {
  return usePostData('/user/list', params)
}

// 获取用户数量
export function userCount(params) {
  return usePostData('/user/count', params)
}

// 获取用户详情
export function userDetail(id) {
  return useGetData('/user/get/' + id)
}

// 启用用户
export function userEnable(id) {
  return usePost('/user/enable/' + id)
}

// 禁用用户
export function userDisable(id) {
  return usePost('/user/disable/' + id)
}

// 删除用户
export function userDelete(id) {
  return usePost('/user/delete/' + id)
}

// 用户统计
export function userStats() {
  return useGetData('/user/stats')
}

// ============ 邀请码管理 API ============

// 获取邀请码列表
export function inviteCodeList(params) {
  return usePostData('/invite-code/list', params)
}

// 获取邀请码数量
export function inviteCodeCount(params) {
  return usePostData('/invite-code/count', params)
}

// 创建邀请码
export function inviteCodeCreate(data) {
  return usePostData('/invite-code/create', data)
}

// 删除邀请码
export function inviteCodeDelete(id) {
  return usePost('/invite-code/delete/' + id)
}