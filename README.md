##Translation REST server.

Simple REST API that might be used to perform dictionary based translation.

Server supports FreeDict TEI file format.
FreeDict TEI dictionary files are not distributed with this project.

You can learn more about the format from this site: http://freedict.org/en/

You can create your own TEI dictionary file or download a needed one from here: https://github.com/freedict/fd-dictionaries

###Make
 > make

###Start server:
 > translate-api [lang from] [lang to] [mode] [path to your dictionary file]
- *mode* Switches translation mode. There are two possible modes: default and prose.

### Example start server:
> ./translate-api nl en prose ~/Desktop/my-nl-to-en.tei

# Translation

###Request parameters:
- *text*    Text to translate
- *from*    Language of the original text 
- *to*      Language of result text
- *maxAlt* Some words might have alternative translations. Specify this param to include certain number of
            words alternative translations into a result text.

###Request example:
> http://localhost:9000/translate?text=%22Goedemorgen%20iedereen%22&from=NL&to=EN&max-alt=2

# Inspection

Debugging endpoint to check what will be the result for particular word or phrase including closest match and distance.

###Request parameters:
- *lang-from* Language of original text
- *lang-to* Language of result text
- *text* Text to inspect

###Request example:
> http://localhost:9000/inspect?text=balsturig&lang-from=NL&lang-to=EN