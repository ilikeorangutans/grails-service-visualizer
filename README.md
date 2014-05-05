grails-service-visualizer
=========================

Visualizes Grails Service and Controller Dependencies using Graphviz

Simple tool I wrote one afternoon to visualize complex Grails apps I was working with. See [my blog post](http://ilikeorangutans.github.io/2014/05/03/using-golang-and-graphviz-to-visualize-complex-grails-applications/) about the topic. 

To Do
-----

This is app is a huge hack and currently can even produce invalid Graphviz files. It's by no means perfect, but allowed me to generate the beautiful graphs reference below:


Sample Images
-------------
![Sample Graph](http://ilikeorangutans.github.io/assets/images/sample-app-01.png)
[App A Full Resolution[PNG 558K]](http://ilikeorangutans.github.io/assets/images/sample-app-01.png)

![Sample Graph](http://ilikeorangutans.github.io/assets/images/sample-app-02.png)
[App B Full Resolution[PNG 393K]](http://ilikeorangutans.github.io/assets/images/sample-app-02.png)


Ideas
-----

- Calculate the transitive number of dependencies a component has and adjust the size of the visual representation based on that. For example, a service used only by a single controller would be fairly small. Another service, used by twenty other services would be much bigger and more prominent. Fairly easy to implement.
- Use afferent and [efferent](http://en.wikipedia.org/wiki/Efferent_coupling) coupling to calculate the [instability](http://en.wikipedia.org/wiki/Software_package_metrics) of services. Originally a metric for packages, this could probably applied to services as well. Could be visualized with border thickness of shapes.
- Support different types of artifacts and give them different shapes. Right now I'm thinking quartz jobs which play a fairly important role in our application. 
- Parse code instead of string searching. Right now I'm simply scanning each source file line by line and search for certain patterns. This works reasonably well, but with some apps I get unwanted values. Sometimes from comments, sometimes other values. If I were to actually parse the source into an AST, all these issues would disappear. It would also allow a host of other interesting things, like: 
- Calculate the [Cyclomatic Complexity](http://en.wikipedia.org/wiki/Cyclomatic_complexity) for each component and use that to colour code the visual representation. Low complexity could be a green tint, high values would be red. 
- Support different languages and frameworks. I was able to implement this tool very easily because Grails' service and controller naming follows an entirely name based convention. That makes finding service definitions trivial. I would love to extend this tool to support different languages and frameworks. 
- Have more ideas? Drop me an email or a comment or send me a pull request! :)