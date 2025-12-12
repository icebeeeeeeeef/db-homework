import axios from "axios";
import router from "next/router";
const instance = axios.create({
    // 这边记得修改你对应的配置文件
    baseURL: "http://localhost:8080",
    withCredentials: true
})

export interface Result<T> {
    code: number,
    msg: string,
    data: T,
}


instance.interceptors.response.use(
    (resp) => {
        if (typeof window !== "undefined") {
            const newToken = resp?.headers?.["x-jwt-token"]
            const newRefreshToken = resp?.headers?.["x-refresh-token"]
            if (newToken) {
                localStorage.setItem("token", newToken)
                // 立刻更新默认请求头，避免导航后第一次请求没带上 token
                instance.defaults.headers.common["Authorization"] = "Bearer " + newToken
            }
            if (newRefreshToken) {
                localStorage.setItem("refresh_token", newRefreshToken)
            }
        }
        return resp
    },
    (err) => {
        console.log(err)
        const status = err?.response?.status
        if (status === 401) {
            window.location.href = "/users/login"
        }
        // 继续抛出，让调用方能拿到错误信息
        return Promise.reject(err)
    }
)

// 在这里让每一个请求都加上 authorization 的头部
instance.interceptors.request.use((req) => {
    if (typeof window !== "undefined") {
        const token = localStorage.getItem("token")
        if (token) {
            // Axios v1 的 headers 可能是 AxiosHeaders，也可能是普通对象，两种方式都兼容
            if (!req.headers) {
                req.headers = {}
            }
            // @ts-ignore
            if (typeof req.headers.set === "function") {
                // @ts-ignore
                req.headers.set("Authorization", "Bearer " + token)
            } else {
                req.headers["Authorization"] = "Bearer " + token
            }
        }
    }
    return req
}, (err) => {
    console.log(err)
    return Promise.reject(err)
})

export default instance
