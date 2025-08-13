
## Frequently Asked Questions

* Can I use URLs in N4L?

Yes, you should enclose them in quotes because they usually include the substring "//", which is also a comment designation.

* Why does it take so long to upload data?

Uploading to a database is a slow process compared to retrieving as there are many checks that have to happen. try to debug your data as far as possible using the text interface in N4L before actually committing to the database.

In addition, *Unicode* decoding is a very slow process so long files seem to take forever to read, never mind the actual database uploading. I don't know of any way to speed this up presently. Unless we know that a file is simple
ASCII encoding, it's easy to get bad character conversion without using this longwinded decoding.

* Why are there relationships that I didn't intend when I browse the data?

Be careful to ensure that you haven't accidentally used any of the annotation markers (e.g. +,-,=) without surroundings spaces in your text, as these will be interpreted as annotations. Use the verbose mode in N4L to debug.

* Why do I see chapters that don't seem to be relevant?

This is probably a result of certain words and phrases belonging to more than one chapter, and thus bridging chapters that you didn't intend. This bridging is intentional, as it allows >"lateral thinking", which is an important source of discovery.