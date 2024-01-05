
const MARKUP_VARIABLES = {
  registry: {},
  nestedElements: [],
  write: function(name, value) {
    this.registry[name] = value;
    this.applyOnDocument();
  },
  read: function(name) {
    return this.registry[name];
  },
  get vars() {
    return Object.keys(this.registry);
  },
  applyOnDocument: function() {

    document.querySelectorAll('*').forEach(element => {
      var child = element.firstChild;
      if(element.innerHTML.includes('%'))
      while(child) {

        const nested = {
          node: child,
          content: child.textContent
        };

        if (
          child.nodeType === Node.TEXT_NODE
          && child.textContent.includes('%')
          && !this.nestedElements.includes(nested)
        ) {
          this.nestedElements.push(nested);
          console.log("Registering Nested-Text-Node: ");
          console.log(nested);
        }

        child = child.nextSibling;

      }
    });

    this.nestedElements.forEach(nested =>
    this.vars.forEach(variable => {
        nested.node.textContent = nested.content.replace(
          '%' + variable + '%',
          this.registry[variable]
        );
    }));

  }
};
