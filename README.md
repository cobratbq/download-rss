download-rss
============

Trivial commandline podcast rss feed downloader just for playing around with the Go programming language. The application currently calls 'wget' (synchronously) to execute the download.

If an argument is provided, that argument is used as the source for downloading the podcasts. If no argument is provided, stdin is used as the source.

Usage:
> download-rss feed.xml

References:
- [Parsing huge xml files with go](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/): Nice example of efficiently parsing xml with go.
