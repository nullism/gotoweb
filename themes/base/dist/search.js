var indexData = null

function loadIndex() {
  return fetch("/index.json")
    .then(response => response.json())
    .then(data => {
      indexData = data
    })
}

function search(query) {

  query = query.toLowerCase()

  console.log("index is", indexData)
  // split query by space
  var parts = query.split(" ")
  console.log("parts are ", parts)
  console.log("query is ", query)

  var hitIds = []
  for (var p of parts) {
    if (indexData.idx[p] !== undefined) {
      hitIds.push(...indexData.idx[p])
    }
  }

  console.log("hitIds are ", hitIds)

  const countMap = hitIds.reduce((acc, val) => {
    acc[val] = (acc[val] || 0) + 1
    return acc
  }, {})

  console.log("countMap is ", countMap)
  var sortedHitIds = Object.keys(countMap).sort(function (a, b) {
    return countMap[b] - countMap[a]
  })
  console.log("newTypesArray is ", sortedHitIds)

  var hits = []
  for (var i = 0; i < sortedHitIds.length; i++) {
    var id = sortedHitIds[i]
    var hit = indexData.docs[id]
    hits.push(hit)
  }
  console.log("hits are ", hits)
  return hits
}

loadIndex().then(() => {
  console.log("search index loaded")
})