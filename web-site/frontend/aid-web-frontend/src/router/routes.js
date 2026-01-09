const routes = [
    // {
    //   path: '/',
    //   redirect: '/page/'
    // },
    {
        path: "/",
        component: () => import("layouts/MainLayout.vue"),
        children: [
            { path: "fileManager", component: () => import("src/pages/FileManagerPage.vue") },
            { path: "file", component: () => import("src/pages/FilePage.vue") },
            { path: "clockingIn", component: () => import("src/pages/ClockingInPage.vue") },
            { path: "timeCalculator", component: () => import("src/pages/TimeCalculatorPage.vue") },
            { path: "rezip", component: () => import("src/pages/RezipPage.vue") },
            { path: "uuid", component: () => import("src/pages/UuidPage.vue") },
            { path: "randString", component: () => import("src/pages/RandStringPage.vue") },
            { path: "codeChange", component: () => import("src/pages/CodeChangePage.vue") },
            { path: "", component: () => import("src/pages/Base64ToImagePage.vue") },
            { path: "encoder", component: () => import("src/pages/EncoderPage.vue") },
            { path: "qrCode", component: () => import("src/pages/QRCodePage.vue") },
        ],
    },

    // Always leave this as last one,
    // but you can also remove it
    {
        path: "/:catchAll(.*)*",
        component: () => import("pages/ErrorNotFound.vue"),
    },
];

export default routes;
