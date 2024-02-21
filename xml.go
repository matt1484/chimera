package chimera

// TODO: I want to support json but I need to be able to handle xml.Name and the use of
// "attr"/"chardata"/"any"/"innerxml" in struct tags. Not sure how prevelant these are
// or even the use of xml as whole since JSON is far more prevalent. Fundamentally I
// just need to follow this guide https://swagger.io/docs/specification/data-models/representing-xml/
// and read all the rules in "encoding/xml".
// TLDR: JSON is way easier, XML TBD
