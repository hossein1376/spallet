const fs = require('fs');

if (process.argv.length !== 4) {
    throw "expected two args"
}

const jsFile = process.argv[2]
const htmlFile = process.argv[3]

try {
    const definition = JSON.stringify(fs.readFileSync(jsFile, 'utf8'))
    const schema = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.29.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.29.0/swagger-ui-bundle.js" crossorigin></script>
<script src="https://unpkg.com/swagger-ui-dist@5.29.0/swagger-ui-standalone-preset.js" crossorigin></script>
<script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        spec: JSON.parse(${definition}),
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
      });
    };
</script>
</body>
</html>`

    fs.writeFileSync(htmlFile, schema)
} catch (err) {
    throw `failed to generate html file: ${err}`
}
