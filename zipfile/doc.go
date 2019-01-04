/*Package zipfile reads & parses zip archives

  outline: zipfile
    zipfile reads & parses zip archives

    functions:
      ZipFile(data)
        opens an archive for reading

    types:
      ZipFile
        a zip archive object
        methods:
          namelist() list
            return a list of files in the archive
          open(filename string) ZipInfo
            open a file for reading
            params:
              filename string
                name of the file in the archive to open
      ZipInfo
        methods:
          read() string
            read the file, returning it's string representation


*/
package zipfile
