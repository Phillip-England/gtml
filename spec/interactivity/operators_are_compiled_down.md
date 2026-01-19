When an operator like "+" or "-" or "*" or any other operation is associated with a signal value, it needs to get compiled down.

Even making a signal equal to something should be compiled down.

In gtml, doing this: `$signal = "Bob"` should result in a dom change in the actual html page. It should change the value of the signal in every place in which it exists within the application.

Same thing if we have something like this: `$signal = $signal + 1`. Imagine `$signal` in that case was `1`. Well, now it would be `2` and it would be visible in all areas of the dom becuase this would compile down into nice javascript which changes the value of the signal thus changing the dom.

This should give us some very simple ways to keep our javascripting abilities and giving us control, while finding a nice middleground where we can produce signals freely where we want with sugar syntax.
