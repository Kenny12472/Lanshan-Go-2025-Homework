// static/js/api.js （重写，覆盖原文件）
const API_BASE = "http://localhost:8080"; // 如需相对路径可改成 ""

// token helpers
function getToken() { return localStorage.getItem("token"); }
function setToken(t) { localStorage.setItem("token", t); }
function removeToken() { localStorage.removeItem("token"); }
function authHeader() {
    const t = getToken();
    return t ? { "Authorization": "Bearer " + t } : {};
}

// 更健壮的 request
async function request(url, method = "GET", data = null, auth = false) {
    try {
        if (!url.startsWith("http")) {
            if (!url.startsWith("/")) url = "/" + url;
            url = API_BASE + url;
        }

        const headers = {};
        if (data != null) headers["Content-Type"] = "application/json";
        if (auth) Object.assign(headers, authHeader());

        const opts = { method, headers };
        if (data != null) opts.body = JSON.stringify(data);

        const res = await fetch(url, opts);
        const text = await res.text();
        const ct = (res.headers.get("content-type") || "").toLowerCase();

        // 专门处理 401：把 token 清掉并跳回登录页（单页应用也可以只清 token）
        if (res.status === 401) {
            try { console.warn("401 from", url, text); } catch (e) { }
            removeToken();
            // 若在前端页面：跳转回登录页（避免循环重定向：仅在页面环境下）
            try { if (typeof window !== 'undefined') window.location.href = "/static/auth.html"; } catch (e) { }
            return null;
        }

        if (!res.ok) {
            console.error(`HTTP ${res.status} ${method} ${url} -> ${text}`);
            // 尝试解析 JSON 错误体
            try { return JSON.parse(text); } catch (e) { return null; }
        }

        if (ct.includes("application/json")) {
            try { return JSON.parse(text); } catch (e) {
                console.error("JSON parse error for", url, "body:", text, e);
                return null;
            }
        }
        return text;
    } catch (e) {
        console.error("request exception:", e, url, method, data);
        return null;
    }
}

// auth
async function apiRegister(username, password) { return await request("/register", "POST", { username, password }); }
async function apiLogin(username, password) {
    const r = await request("/login", "POST", { username, password });
    if (r && r.token) setToken(r.token);
    return r;
}

// articles
async function apiCreateArticle(title, content, status) { return await request("/articles", "POST", { title, content, status }, true); }
async function apiListArticles(page = 1, page_size = 10, query = "", followedFirst = false) {
    const q = `?page=${page}&page_size=${page_size}` + (query ? `&query=${encodeURIComponent(query)}` : "") + (followedFirst ? `&followed_first=1` : "");
    return await request("/articles" + q, "GET");
}
async function apiGetArticle(id) { return await request("/articles/" + id, "GET"); }
async function apiUpdateArticle(id, title, content, status) { return await request("/articles/" + id, "PUT", { title, content, status }, true); }
async function apiDeleteArticle(id) { return await request("/articles/" + id, "DELETE", null, true); }

// comments
async function apiPostComment(articleId, content) { return await request(`/articles/${articleId}/comments`, "POST", { content }, true); }
async function apiListComments(articleId) { return await request(`/articles/${articleId}/comments`, "GET"); }

// likes (toggle)
async function apiToggleArticleLike(articleId) { return await request(`/articles/${articleId}/like`, "POST", null, true); }
async function apiToggleCommentLike(commentId) { return await request(`/comments/${commentId}/like`, "POST", null, true); }

// follows
async function apiToggleFollow(userId) { return await request(`/users/${userId}/follow`, "POST", null, true); }
async function apiGetMyFollows() { return await request("/me/follows", "GET", null, true); }

// profile
async function apiGetMyProfile() { return await request("/me/profile", "GET", null, true); }
async function apiUpdateMyProfile(display_name, bio) { return await request("/me/profile", "PUT", { display_name, bio }, true); }
async function apiGetUserProfile(userId) { return await request("/users/" + userId + "/profile", "GET"); }