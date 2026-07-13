//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

extern void handleTrayAction(int);

static NSMenuItem* createItem(const char* title, int tag) {
	NSString* str = [NSString stringWithUTF8String:title];
	NSMenuItem* item = [[NSMenuItem alloc] initWithTitle:str action:NULL keyEquivalent:@""];
	[item setTag:tag];
	return [item autorelease];
}

void start_tray(const unsigned char* data, int len, const char* tooltip) {
	// dispatch_async 到主线程，因为 CGO 从 Go 协程调用时不在主线程，
	// 而 NSStatusBar / NSWindow 必须在主线程创建。
	dispatch_async(dispatch_get_main_queue(), ^{
		@autoreleasepool {
			NSStatusItem* item = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
			[item retain];

			// 设置图标
			NSData* imgData = [NSData dataWithBytes:(void*)data length:len];
			NSImage* img = [[NSImage alloc] initWithData:imgData];
			[img setSize:NSMakeSize(18, 18)];
			item.button.image = img;
			[img release];

			// 设置工具提示
			item.button.toolTip = [NSString stringWithUTF8String:tooltip];

			// 创建菜单
			NSMenu* menu = [[NSMenu alloc] init];
			[menu addItem:createItem("显示窗口", 1)];
			[menu addItem:createItem("隐藏窗口", 2)];
			[menu addItem:[NSMenuItem separatorItem]];
			[menu addItem:createItem("退出", 3)];

			// 统一设置菜单项 action
			SEL sel = NSSelectorFromString(@"trayAction:");
			if (![NSApp respondsToSelector:sel]) {
				IMP imp = imp_implementationWithBlock(^(id self_, id sender) {
					int tag = (int)[(NSMenuItem*)sender tag];
					handleTrayAction(tag);
				});
				class_addMethod([NSApp class], sel, imp, "v@:@");
			}

			for (NSMenuItem* mi in [menu itemArray]) {
				if (![mi isSeparatorItem]) {
					[mi setTarget:NSApp];
					[mi setAction:sel];
				}
			}

			item.menu = menu;
		}
	});
}
*/
import "C"