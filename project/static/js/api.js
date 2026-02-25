const API_BASE = "http://localhost:8080";

// ��ѧ��ע�ͣ� token helpers
function getToken() { return localStorage.getItem("token"); }
function setToken(t) { localStorage.setItem("token", t); }
function removeToken() { localStorage.removeItem("token"); }
function authHeader() {
    const t = getToken();
    return t ? { "Authorization": "Bearer " + t } : {};
}

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

        if (!res.ok) {
            console.error(`HTTP ${res.status} ${method} ${url}`, text);
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

async function apiRegister(username, password) { return await request("/register", "POST", { username, password }); }
async function apiLogin(username, password) {
    const r = await request("/login", "POST", { username, password });
    if (r && r.token) setToken(r.token);
    return r;
}

async function apiCreateArticle(title, content, status) { return await request("/articles", "POST", { title, content, status }, true); }
async function apiListArticles(page = 1, page_size = 10, query = "", followedFirst = false) {
    const q = `?page=${page}&page_size=${page_size}` + (query ? `&query=${encodeURIComponent(query)}` : "") + (followedFirst ? `&followed_first=1` : "");
    return await request("/articles" + q, "GET");
}
async function apiGetArticle(id) { return await request("/articles/" + id, "GET"); }
async function apiUpdateArticle(id, title, content, status) { return await request("/articles/" + id, "PUT", { title, content, status }, true); }
async function apiDeleteArticle(id) { return await request("/articles/" + id, "DELETE", null, true); }

async function apiPostComment(articleId, content) { return await request(`/articles/${articleId}/comments`, "POST", { content }, true); }
async function apiListComments(articleId) { return await request(`/articles/${articleId}/comments`, "GET"); }

async function apiToggleArticleLike(articleId) { return await request(`/articles/${articleId}/like`, "POST", null, true); }
async function apiToggleCommentLike(commentId) { return await request(`/comments/${commentId}/like`, "POST", null, true); }

async function apiToggleFollow(userId) { return await request(`/users/${userId}/follow`, "POST", null, true); }
async function apiGetMyFollows() { return await request("/me/follows", "GET", null, true); }


async function apiGetMyProfile() { return await request("/me/profile", "GET", null, true); }
async function apiUpdateMyProfile(display_name, bio) { return await request("/me/profile", "PUT", { display_name, bio }, true); }
async function apiGetUserProfile(userId) { return await request("/users/" + userId + "/profile", "GET"); }


