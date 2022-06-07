package goscript

// this package will contain the preprocessor for goscript

/* generateFQSC will apply preprocessing steps to the applications source code
   The following steps will be applied:
	- delete import directives
	- delete application directive
	- delete external directives
	- prefix all symbols with the hash of their module (a function getJWT from the package jwt would become hash_getJWT)
	- replace external symbol references with their expected name (jwt.getJWT becomes hash_getJWT)
		- dots can be contained in property access, member access, numbers and strings
		- numbers shouldnt be a problem since we can just ignore them based on required property names
		- we need to detect all strings, remember their starting and ending indices in the string and then bounds check wether or not the
		  match we have is inside a string zone
	- merge all source files
*/
func generateFQSC(source *ApplicationSource) (string, error) {
	return "", nil
}
