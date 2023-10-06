export default defineNuxtRouteMiddleware((to, from) => {
    if (!process.server) {
        const token = useCookie('token')
        const publicPages = ['/authorization', '/registration', '/'];
        const beforeAuthPages = ['/authorization', '/registration']
        const authRequired = !publicPages.includes(to.path);
        const beforeAuth = beforeAuthPages.includes(to.path);
        if (authRequired && !token.value) {
            return navigateTo('/authorization')
        }
        if (beforeAuth && token.value) {
            return navigateTo('/')
        }
    }
})

