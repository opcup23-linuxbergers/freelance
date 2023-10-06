// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    ssr: false,
    css: ["@/assets/scss/main.scss"],
    runtimeConfig: {
        public: {
            apiBase: 'https://opcup23-api.carried.ru',
            domainName: 'https://opcup23-api.carried.ru'
        }
    }
})
