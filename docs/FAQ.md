
## Frequently Asked Questions

* Can I use URLs in N4L?

Yes, you should enclose them in quotes because they usually include the substring "//", which is also a comment designation.

* Why does it take so long to upload data?

Uploading to a database is a slow process compared to retrieving as there are many checks that have to happen. try to debug your data as far as possible using the text interface in N4L before actually committing to the database.

* Why are there relationships that I didn't intend when I browse the data?

Be careful to ensure that you haven't accidentally used any of the annotation markers (e.g. +,-,=) without surroundings spaces in your text, as these will be interpreted as annotations. Use the verbose mode in N4L to debug.