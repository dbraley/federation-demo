async function shippingEstimate({parents, args}) {
    console.log(parents)
    console.log(args)
    return parents.map(p => parseInt(p.upc) + 1000)
}

self.addMultiParentGraphQLResolvers({
    "Product.shippingEstimate": shippingEstimate
})
