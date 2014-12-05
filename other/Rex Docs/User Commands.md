
User Commands

============================================================
Loaded by an init script in the "Libs" addon:

Generate a place holder template.
	
	rubble:placeholder name
	
	Creates a template named <name> that returns a simple placeholder message,
	may be overridden by later definitions.

Makes Rubble skip a file.
	
	rubble:skipfile name
	
	name is the file's BASE NAME not it's path!
	Returns unchanged.
	
Makes Rubble not write a file out after generation..
	
	rubble:nowritefile name
	
	name is the file's BASE NAME not it's path!
	Returns unchanged.

Marks a file as containing creature graphics.
	
	rubble:graphicsfile name
	
	name is the file's BASE NAME not it's path!
	Returns unchanged.

Marks a file as a prep file.
	
	rubble:prepfile name
	
	name is the file's BASE NAME not it's path!
	Returns unchanged.

============================================================
Loaded by a loader script in the "Libs" addon:

Checks if a named addon is active.
	
	rubble:addonactive name
	
	name needs to be the FULL name of the addon! (for example "Libs/Base")
	The name check is case-sensitive!
	Returns true or false.

Conditionally activates an addon from a load script.
	
	rubble:activateaddon me addon
	
	me and addon both need to be the full addon names! 
	me should be the addon doing the checking, addon should be the one checked.
	If me is not active then addon is not activated.
	If this is called from anything other than a loader script nothing will happen!

Rubble version check, aborts if the current version is not the same as or newer than the requested version.
	
	rubble:checkversion addon version
	
	If the current version is newer a message is printed saying so.
	addon is used as the name of what is requesting a specific version and does not have to be a valid addon name.
