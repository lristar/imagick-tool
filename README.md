# Creating pdf preview by making an image from the first page with Imagick

Build a docker image with the following command 

    docker build -t Imagick.

Then create a new docker container with this command

    docker run -p 900:900 imagick

Now you have a docker container up&running so you can make a request to make an image 

    curl -F "file=@test1234.pdf" localhost:900/convert --output test1234.jpg

You can replace your pdf file instead of ``test1234.pdf`` so you will have output in your current directory name is ``test1234.jpg``

If there is any issue or improvement please contribute and make it better

This is just to test how Imagick works in Go and not suitable for production

