// Code generated by gtml; DO NOT EDIT.
// +build ignore
// v0.1.0 | you may see errors with types, you'll need to manage your own imports
// type support coming soon!

package main

import "strings"

func gtmlFor[T any](slice []T, callback func(i int, item T) string) string {
var builder strings.Builder
for i, item := range slice {
	builder.WriteString(callback(i, item))
}
return builder.String()
}

func gtmlIf(condition bool, fn func() string) string {
if condition {
	return fn()
}
return ""
}

func gtmlElse(condition bool, fn func() string) string {
if !condition {
	return fn()
}
return ""
}

func gtmlSlot(contentFunc func() string) string {
return contentFunc()
}

func DiningMenu(foodThree string) string {
	diningmenu := func() string {
		var diningmenuBuilder strings.Builder
		diningmenuBuilder.WriteString(`<div _component="DiningMenu" _id="0"><h1>Welcome!</h1><p>Please take a look at our menu, ask if you have questions!</p><foodlist food-one="Pizza" food-two="Tacos" food-three="`)
		diningmenuBuilder.WriteString(foodThree)
		diningmenuBuilder.WriteString(`"></foodlist></div>`)
		return diningmenuBuilder.String()
	}
	return diningmenu()
}

func IfElement(isLoggedIn bool) string {
	ifelement := func() string {
		var ifelementBuilder strings.Builder
		isLoggedInElse1 := gtmlElse(isLoggedIn, func() string {
			var isLoggedInBuilder strings.Builder
			isLoggedInBuilder.WriteString(`<div _else="isLoggedIn" _id="1"><p>you are not logged in!</p></div>`)
			if !isLoggedIn {
					return isLoggedInBuilder.String()
			}
			return ""
		})
		ifelementBuilder.WriteString(`<div _component="IfElement" _id="0">`)
		ifelementBuilder.WriteString(isLoggedInElse1)
		ifelementBuilder.WriteString(`</div>`)
		return ifelementBuilder.String()
	}
	return ifelement()
}

func ForCustomSlice(Guests []Guest) string {
	forcustomslice := func() string {
		var forcustomsliceBuilder strings.Builder
		guestFor1 := gtmlFor(Guests, func(i int, guest Guest) string {
			var guestBuilder strings.Builder
			guestBuilder.WriteString(`<ul _for="guest of Guests []Guest" _id="1"><p>`)
			guestBuilder.WriteString(guest.Name)
			guestBuilder.WriteString(`</p></ul>`)
			return guestBuilder.String()
		})
		forcustomsliceBuilder.WriteString(`<div _component="ForCustomSlice" _id="0">`)
		forcustomsliceBuilder.WriteString(guestFor1)
		forcustomsliceBuilder.WriteString(`</div>`)
		return forcustomsliceBuilder.String()
	}
	return forcustomslice()
}

func ForStringSlice(colors []string) string {
	forstringslice := func() string {
		var forstringsliceBuilder strings.Builder
		colorFor1 := gtmlFor(colors, func(i int, color string) string {
			var colorBuilder strings.Builder
			colorBuilder.WriteString(`<ul _for="color of colors []string" _id="1"><p>`)
			colorBuilder.WriteString(color)
			colorBuilder.WriteString(`</p></ul>`)
			return colorBuilder.String()
		})
		forstringsliceBuilder.WriteString(`<div _component="ForStringSlice" _id="0">`)
		forstringsliceBuilder.WriteString(colorFor1)
		forstringsliceBuilder.WriteString(`</div>`)
		return forstringsliceBuilder.String()
	}
	return forstringslice()
}

func GreetingCard(name string, guestFirstName string, colors []string) string {
	greetingcard := func() string {
		var greetingcardBuilder strings.Builder
		greetingslotPlaceholder1 := func() string {
			messageSlot2 := gtmlSlot(func() string {
				var messageBuilder strings.Builder
				messageBuilder.WriteString(`<div _slot="message" _id="2"><p>testin!</p></div>`)
				return messageBuilder.String()
			})
			loopSlot3 := gtmlSlot(func() string {
				var loopBuilder strings.Builder
				colorFor4 := gtmlFor(colors, func(i int, color string) string {
					var colorBuilder strings.Builder
					colorBuilder.WriteString(`<ul _for="color of colors []string" _id="4"><li>`)
					colorBuilder.WriteString(color)
					colorBuilder.WriteString(`</li></ul>`)
					return colorBuilder.String()
				})
				loopBuilder.WriteString(`<div _slot="loop" _id="3">`)
				loopBuilder.WriteString(colorFor4)
				loopBuilder.WriteString(`</div>`)
				return loopBuilder.String()
			})
			return GreetingSlot(messageSlot2, guestFirstName, "20", loopSlot3)
		}
		greetingcardBuilder.WriteString(`<div _component="GreetingCard" _id="0"><h1>`)
		greetingcardBuilder.WriteString(name)
		greetingcardBuilder.WriteString(`</h1>`)
		greetingcardBuilder.WriteString(greetingslotPlaceholder1())
		greetingcardBuilder.WriteString(`</div>`)
		return greetingcardBuilder.String()
	}
	return greetingcard()
}

func GreetingSlot(message string, name string, age string, loop string) string {
	greetingslot := func() string {
		var greetingslotBuilder strings.Builder
		greetingslotBuilder.WriteString(`<div _component="GreetingSlot" _id="0">`)
		greetingslotBuilder.WriteString(message)
		greetingslotBuilder.WriteString(`<h1>Hello, `)
		greetingslotBuilder.WriteString(name)
		greetingslotBuilder.WriteString(`</h1> <p>you are `)
		greetingslotBuilder.WriteString(age)
		greetingslotBuilder.WriteString(` years old!</p>`)
		greetingslotBuilder.WriteString(loop)
		greetingslotBuilder.WriteString(`</div>`)
		return greetingslotBuilder.String()
	}
	return greetingslot()
}

func GuestMesh(someTitle string, guests []Guest, isAdmin bool, colors []string, loggedIn bool) string {
	guestmesh := func() string {
		var guestmeshBuilder strings.Builder
		guestFor1 := gtmlFor(guests, func(i int, guest Guest) string {
			var guestBuilder strings.Builder
			itemFor2 := gtmlFor(guest.Items, func(i int, item Item) string {
				var itemBuilder strings.Builder
				colorFor3 := gtmlFor(item.Colors, func(i int, color Color) string {
					var colorBuilder strings.Builder
					isAdminIf4 := gtmlIf(isAdmin, func() string {
						var isAdminBuilder strings.Builder
						isAdminBuilder.WriteString(`<div _if="isAdmin" _id="4"><p>Hello Admin!</p></div>`)
						if isAdmin {
							return isAdminBuilder.String()
						}
						return ""
					})
					colorBuilder.WriteString(`<div _for="color of item.Colors []Color" _id="3"><p>`)
					colorBuilder.WriteString(color.Hue)
					colorBuilder.WriteString(`</p><p>`)
					colorBuilder.WriteString(color.Name)
					colorBuilder.WriteString(`</p>`)
					colorBuilder.WriteString(isAdminIf4)
					colorBuilder.WriteString(`</div>`)
					return colorBuilder.String()
				})
				itemBuilder.WriteString(`<div _for="item of guest.Items []Item" _id="2"><p>`)
				itemBuilder.WriteString(item.Name)
				itemBuilder.WriteString(`</p><p>`)
				itemBuilder.WriteString(item.Price)
				itemBuilder.WriteString(`</p>`)
				itemBuilder.WriteString(colorFor3)
				itemBuilder.WriteString(`</div>`)
				return itemBuilder.String()
			})
			guestBuilder.WriteString(`<div _for="guest of guests []Guest" _id="1"><h1>`)
			guestBuilder.WriteString(guest.Name)
			guestBuilder.WriteString(`</h1><p>The guest has brought the following items:</p>`)
			guestBuilder.WriteString(itemFor2)
			guestBuilder.WriteString(`</div>`)
			return guestBuilder.String()
		})
		colorFor5 := gtmlFor(colors, func(i int, color string) string {
			var colorBuilder strings.Builder
			colorBuilder.WriteString(`<div _for="color of colors []string" _id="5"><p>`)
			colorBuilder.WriteString(color)
			colorBuilder.WriteString(`</p><p>`)
			colorBuilder.WriteString(color)
			colorBuilder.WriteString(`</p></div>`)
			return colorBuilder.String()
		})
		loggedInIf6 := gtmlIf(loggedIn, func() string {
			var loggedInBuilder strings.Builder
			loggedInBuilder.WriteString(`<div _if="loggedIn" _id="6"><p>Logged in!</p></div>`)
			if loggedIn {
				return loggedInBuilder.String()
			}
			return ""
		})
		guestmeshBuilder.WriteString(`<div _component="GuestMesh" _id="0"><h1>`)
		guestmeshBuilder.WriteString(someTitle)
		guestmeshBuilder.WriteString(`</h1>`)
		guestmeshBuilder.WriteString(guestFor1)
		guestmeshBuilder.WriteString(colorFor5)
		guestmeshBuilder.WriteString(loggedInIf6)
		guestmeshBuilder.WriteString(`</div>`)
		return guestmeshBuilder.String()
	}
	return guestmesh()
}

func IfElement(isLoggedIn bool) string {
	ifelement := func() string {
		var ifelementBuilder strings.Builder
		isLoggedInIf1 := gtmlIf(isLoggedIn, func() string {
			var isLoggedInBuilder strings.Builder
			isLoggedInBuilder.WriteString(`<div _if="isLoggedIn" _id="1"><p>you are logged in!</p></div>`)
			if isLoggedIn {
				return isLoggedInBuilder.String()
			}
			return ""
		})
		ifelementBuilder.WriteString(`<div _component="IfElement" _id="0">`)
		ifelementBuilder.WriteString(isLoggedInIf1)
		ifelementBuilder.WriteString(`</div>`)
		return ifelementBuilder.String()
	}
	return ifelement()
}

func PlaceholderBasic() string {
	buttonplaceholderPlaceholder0 := func() string {
		return ButtonPlaceholder()
	}
	return buttonplaceholderPlaceholder0()
}

func ButtonPlaceholder() string {
	buttonplaceholder := func() string {
		var buttonplaceholderBuilder strings.Builder
		buttonplaceholderBuilder.WriteString(`<div _component="ButtonPlaceholder" _id="0"><button>Submit</button></div>`)
		return buttonplaceholderBuilder.String()
	}
	return buttonplaceholder()
}

func PlaceholderPropInAttr(name string) string {
	placeholderpropinattr := func() string {
		var placeholderpropinattrBuilder strings.Builder
		greetingPlaceholder1 := func() string {
			return Greeting(name)
		}
		placeholderpropinattrBuilder.WriteString(`<div _component="PlaceholderPropInAttr" _id="0">`)
		placeholderpropinattrBuilder.WriteString(greetingPlaceholder1())
		placeholderpropinattrBuilder.WriteString(`</div>`)
		return placeholderpropinattrBuilder.String()
	}
	return placeholderpropinattr()
}

func Greeting(name string) string {
	greeting := func() string {
		var greetingBuilder strings.Builder
		greetingBuilder.WriteString(`<div _component="Greeting" _id="0"><h1>`)
		greetingBuilder.WriteString(name)
		greetingBuilder.WriteString(`</h1></div>`)
		return greetingBuilder.String()
	}
	return greeting()
}

func NameTag(firstName string, message string) string {
	nametag := func() string {
		var nametagBuilder strings.Builder
		nametagBuilder.WriteString(`<div _component="NameTag" _id="0"><h1>`)
		nametagBuilder.WriteString(firstName)
		nametagBuilder.WriteString(`</h1><p>`)
		nametagBuilder.WriteString(message)
		nametagBuilder.WriteString(`</p></div>`)
		return nametagBuilder.String()
	}
	return nametag()
}

func PlaceholderWithAttrs() string {
	nametagPlaceholder0 := func() string {
		return NameTag("Melody", "you are amazing!")
	}
	return nametagPlaceholder0()
}

func RunePipe(age string) string {
	runepipe := func() string {
		var runepipeBuilder strings.Builder
		greetingPlaceholder1 := func() string {
			return Greeting(age)
		}
		runepipeBuilder.WriteString(`<div _component="RunePipe" _id="0"><p>Sally is `)
		runepipeBuilder.WriteString(age)
		runepipeBuilder.WriteString(` years old</p>`)
		runepipeBuilder.WriteString(greetingPlaceholder1())
		runepipeBuilder.WriteString(`</div>`)
		return runepipeBuilder.String()
	}
	return runepipe()
}

func Greeting(age string) string {
	greeting := func() string {
		var greetingBuilder strings.Builder
		greetingBuilder.WriteString(`<div _component="Greeting" _id="0"><h1>This age was piped in!</h1> <p>`)
		greetingBuilder.WriteString(age)
		greetingBuilder.WriteString(`</p></div>`)
		return greetingBuilder.String()
	}
	return greeting()
}

func RuneProp(name string) string {
	runeprop := func() string {
		var runepropBuilder strings.Builder
		runepropBuilder.WriteString(`<div _component="RuneProp" _id="0"><p>Hello, `)
		runepropBuilder.WriteString(name)
		runepropBuilder.WriteString(`!</p></div>`)
		return runepropBuilder.String()
	}
	return runeprop()
}

func RuneAttrProp(name string) string {
	runeattrprop := func() string {
		var runeattrpropBuilder strings.Builder
		runeattrpropBuilder.WriteString(`<div _component="RuneAttrProp" _id="0"><p class="text-sm `)
		runeattrpropBuilder.WriteString(name)
		runeattrpropBuilder.WriteString(`">My class is set to what?</p></div>`)
		return runeattrpropBuilder.String()
	}
	return runeattrprop()
}

func RuneSlot(top string, bottom string) string {
	runeslot := func() string {
		var runeslotBuilder strings.Builder
		runeslotBuilder.WriteString(`<div _component="RuneSlot" _id="0">`)
		runeslotBuilder.WriteString(top)
		runeslotBuilder.WriteString(`<h1>🥪</h1>`)
		runeslotBuilder.WriteString(bottom)
		runeslotBuilder.WriteString(`</div>`)
		return runeslotBuilder.String()
	}
	return runeslot()
}

func Sandwich() string {
	runeslotPlaceholder0 := func() string {
		bottomSlot1 := gtmlSlot(func() string {
			var bottomBuilder strings.Builder
			bottomBuilder.WriteString(`<p _slot="bottom" _id="1">I am on bottom</p>`)
			return bottomBuilder.String()
		})
		topSlot2 := gtmlSlot(func() string {
			var topBuilder strings.Builder
			topBuilder.WriteString(`<p _slot="top" _id="2">I am on top</p>`)
			return topBuilder.String()
		})
		return RuneSlot(topSlot2, bottomSlot1)
	}
	return runeslotPlaceholder0()
}

func RuneVal(colors []string) string {
	runeval := func() string {
		var runevalBuilder strings.Builder
		colorFor1 := gtmlFor(colors, func(i int, color string) string {
			var colorBuilder strings.Builder
			colorBuilder.WriteString(`<ul _for="color of colors []string" _id="1"><p>`)
			colorBuilder.WriteString(color)
			colorBuilder.WriteString(`</p></ul>`)
			return colorBuilder.String()
		})
		runevalBuilder.WriteString(`<div _component="RuneVal" _id="0">`)
		runevalBuilder.WriteString(colorFor1)
		runevalBuilder.WriteString(`</div>`)
		return runevalBuilder.String()
	}
	return runeval()
}

func Echo(message string) string {
	echo := func() string {
		var echoBuilder strings.Builder
		echoBuilder.WriteString(`<div _component="Echo" _id="0"><p>`)
		echoBuilder.WriteString(message)
		echoBuilder.WriteString(`</p><p>`)
		echoBuilder.WriteString(message)
		echoBuilder.WriteString(`</p></div>`)
		return echoBuilder.String()
	}
	return echo()
}

