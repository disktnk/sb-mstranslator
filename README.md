# Microsoft Translator API Plugin

API detail, see: https://www.microsoft.com/en-us/translator/translatorapi.aspx

The API needs "Client ID" and "Client Secret", see: https://datamarket.azure.com/developer/applications

## Usage

```sql
CREATE STATE token TYPE mstranslate WITH
    client_id = "<Client ID>",
    client_secret = "<Client Secret>";
```

```sql
EVAL mstranslate("token", "ja", "en", "美しい日本語");
```

Then, will get "Beautiful Japanese".
