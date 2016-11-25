##TranslationAPI REST server.

Simple REST API that can be used to perform translation.

For now only FreeDict TEI dictionary files compatible format is supported.
FreeDict TEI dictionary files are not distributed with this project.
You can download files for many language pairs and learn more about the format from this site http://freedict.org/en/

###Start server:
 > go run server.go [lang from] [lang to] [path to your dictionary file]

###Request parameters:
- *text*    Text to translate
- *from*    Language of the original text 
- *to*      Language of result text
- *max-alt* Some words might have alternative translations. Specify this param to include certain number of
            words alternative translations into a result text.

###Request sample:
> http://localhost:9000/translate?text=%22Goedemorgen%20iedereen%22&from=NL&to=EN&max-alt=2
