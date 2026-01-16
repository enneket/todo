"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var vue_1 = require("vue");
var pinia_1 = require("pinia");
var vue_i18n_1 = require("vue-i18n");
require("./style.css");
var App_vue_1 = require("./App.vue");
var en_json_1 = require("./locales/en.json");
var zh_json_1 = require("./locales/zh.json");
var i18n = (0, vue_i18n_1.createI18n)({
    legacy: false,
    locale: 'zh',
    fallbackLocale: 'en',
    messages: {
        en: en_json_1.default,
        zh: zh_json_1.default
    }
});
var pinia = (0, pinia_1.createPinia)();
var app = (0, vue_1.createApp)(App_vue_1.default);
app.use(pinia);
app.use(i18n);
app.mount('#app');
