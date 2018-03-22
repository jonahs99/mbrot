# Yet another mandelbrot set visualization

![sample output](example/out.png)

Go has built-in support for complex numbers! The visualization above is colored according to the number of iterations before the coordinate diverges. Multiple points are randomly sampled within each pixel to get some anti-aliasing.

It will perform [sobel](https://github.com/jonahs99/sobel) edge detection with the -sobel option:

![sample with sobel edge detection](example/sobel.png)