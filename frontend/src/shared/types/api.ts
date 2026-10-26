export interface ApiResult<T = unknown> {
  code: number
  msg: string
  data: T
  trace_id?: string
}

export interface PageResult<T> {
  list: T[]
  total: number
}

export interface PageInfo {
  page_id: number
  page_level: number
  page_name: string
  page_path: string
  page_type: number
  button_code?: string
  api_method?: string
  api_path?: string
  child_page: PageInfo[]
  buttons?: PageInfo[]
}

export interface AccountInfo {
  account_id: number
  role_id: number
  role_name: string
  page: PageInfo[]
}

export interface LoginResult {
  id: number
  token: string
  user_name: string
  mobile: string
  account_number: string
  status: number
  last_login_time: number
  create_time: number
}

export interface Account {
  id: number
  user_name: string
  mobile: string
  account_number: string
  password?: string
  status: number
  last_login_time: number
  create_time: number
  update_time: number
}

export interface Role {
  id: number
  role_name: string
  create_time: number
}

export interface MenuPage {
  id: number
  page_name: string
  url: string
  parent_id: number
  page_level: number
  page_type: number
  button_code: string
  api_method: string
  api_path: string
  create_time: number
  parent_name?: string
}

export interface LoginLog {
  id: number
  account_id: number
  account_number: string
  login_ip: string
  user_agent: string
  status: number
  message: string
  create_time: number
}
