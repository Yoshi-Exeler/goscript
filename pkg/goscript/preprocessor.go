package goscript

// this package will contain the preprocessor for goscript

/* PreprocessAppSource will apply preprocessing steps to the applications source code
   The following steps will be applied:
	- delete all comments
	- delete import directives
	- delete application directive
	- delete external directives
	- prefix all symbols with the hash of their module (a function getJWT from the package jwt would become hash_getJWT)
	- replace external symbol references with their expected name (jwt.getJWT becomes hash_getJWT)
	- merge all source files
*/
func PreprocessAppSource(source *ApplicationSource) string { return "" }
