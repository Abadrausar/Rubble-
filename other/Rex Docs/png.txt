PACKAGE DOCUMENTATION

package png
    import "dctech/rex/commands/png"



FUNCTIONS

func Command_Load(script *rex.Script, params []*rex.Value)
    Load a PNG image from a "string".

	png:load image_data

    If loading from a file be careful not to mess up the image data. Returns
    a set of nested indexables ([image x y]) pixels are represented by 32
    bit alpha premultiplied colors.


