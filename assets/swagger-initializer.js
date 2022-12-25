// Script contains Swagger UI configuration and initialization.
window.onload = function() {
  // Begin Swagger UI call region (version for configuration with URL query params).
  window.ui = SwaggerUIBundle({
    dom_id: '#swagger-ui',
    queryConfigEnabled: true,
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
  // End Swagger UI call region.
};
