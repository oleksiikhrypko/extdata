window.onload = function () {
    //<editor-fold desc="Changeable Configuration Block">
    window.ui = SwaggerUIBundle({
        configUrl: '/docs/config.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout",
        validatorUrl: "none",
        showMutatedRequest: false,
        requestSnippetsEnabled: true,
        requestSnippets: {
            generators: {
                curl_bash: {
                    title: "cURL (bash)",
                    syntax: "bash"
                },
            },
            defaultExpanded: true,
            languages: ["curl_bash"],
        },
        requestInterceptor: function (request) {
            return request
        }
    });
    //</editor-fold>
};
