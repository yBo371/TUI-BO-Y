from textual.app import App, ComposeResult
from textual.widgets import Button, Static


class MyApp(App):
    def compose(self) -> ComposeResult:
        yield Static("欢迎使用 BO-Y 终端面板")
        yield Button("点击我", id="hello")
        yield Button("退出", id="quit")

    def on_button_pressed(self, event: Button.Pressed) -> None:
        if event.button.id == "hello":
            self.query_one(Static).update("你点击了按钮")

        if event.button.id == "quit":
            self.exit()


if __name__ == "__main__":
    MyApp().run()