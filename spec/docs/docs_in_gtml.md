# Our Documentation is Made with GTML

## We Should Test GTML By Making A Documentation Website
What better way to ensure that GTML is functions as expected than by creating the documentation site with it?

## Located in  ./www
The gtml project where the doucmentation exists will be located at `./www` at this projects root (not reliative to this dir). This file will be a standard gtml project as outlined by this spec.


## Build With Preinstalled Components
The documentation website ought to be built with preinstalled gtml components such that they showcase how the components look and feel when in use. If you find yourself building out the documentation website, and you need to create more components to get the job done, please create them in the spec, add them to the actual code for the project, and dont forget to test them please. Also, let me know if you are doing this. The project may grow overtime, and we may add to the documentation. If you need new components for the documentation site, just tell me and then write a spec for them, write the go code to place them in preinstalled components upon init, and then write tests for them.

## Links Within the Site
Inside a gtml site, the links do not need to point to a `html` like this:
```html
<a href='/docs.html'></a>
```

We are free to do:
```html
<a href='/docs'></a>
```
