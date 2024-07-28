Implement a very simple full-text search engine on the included products_small.json available [here](./product_small.json).

The goal of a full-text search engine is to optimize lookups of items by string tokens, often at the expense of pre-processing the data and using lots of memory (i.e. the same data stored in multiple different ways).

If you are an experienced search engineer, feel free to build on your existing domain knowledge, though none is expected for this exercise.

Requirements and stretch goals:

Implement a simple and fast way to find all product names that match keyboard (note that we are only asking to match the product names and not other fields, such as material or description)
Implement AND conditions, e.g. steel keyboard should match only item names which contain both steel and keyboard somewhere (not necessarily in that order)
Implement OR conditions, e.g. steel keyboard should match item names steel table and wooden keyboard because they each contain one of our search terms
