function shippingEstimate(object) {
      // free for expensive items
      if (object.price > 1000) return 0;
      // estimate is based on weight
      return object.weight * 0.5;
}

async function shippingEstimateResolver({parents}) {
    return parents.map(p => shippingEstimate(p));
}

self.addMultiParentGraphQLResolvers({
    "Product.shippingEstimate": shippingEstimateResolver
})
