# Türkçe Wordle Oyunu

## Açıklama

Wordle oyunu altı denemede 5 harfli bir kelimeyi tahmin etme oyunudur.
Oyun her tahmin edilen kelime için hedef kelime içerisinde mevcut ama yanlış konumdaki ya da hedef kelime içerisinde mevcut ve doğru konumdaki harfleri farklı renklerde geri bildirimde bulunur.

`hedef kelime içerisinde mevcut ama yanlış konumdaki harf => Sarı`

`hedef kelime içerisinde mevcut ve doğru konumdaki harf => Yeşil`

Eğer tüm harfler hedef kelime içerisinde mevcut ve doğru konumda ise tüm harfler yeşil olur ve oyunu kazanırsınız.
Eğer 6 denemede hedef kelimeyi tahmin edemezseniz oyunu kaybedersiniz.

## Ekran Görüntüsü

![Türkçe Wordle](./wordle.png)

## Gereksinimler

- [Go 1.23.1 ya da üstü](https://go.dev/)

## Kullanılan Kütüphaneler

- [Ebitengine V2](https://github.com/hajimehoshi/ebiten) - Game Engine
- [etxt](https://github.com/tinne26/etxt) - Text Rendering

## Nasıl Çalıştırılır?

```bash
git clone https://github.com/DTVegaArchChapter/GameProgramming.git
cd GameProgramming/2d-games/wordle/golang
go run main.go
```

## Neler Öğrendik?

### Game Loop

Ebitengine'de `game loop` inşa etmek için bir `struct` tipi yanımlayıp `ebiten.Game` interface'ini implemente ettikten sonra nesneyi yaratıp `ebiten.RunGame()` fonksiyonuna vermek gerekiyor.

`ebiten.Game` interface'ini implemente etmek için `Update`, `Draw` ve `Layout` fonksiyonlarını tanımlamak gerekmektedir.

- `Layout` fonksiyonu window yeniden boyutlandırıldığında oyunun nasıl ölçekleneceğini belirlemenize izin veren fonksiyondur.
- `Update` fonksiyonu oyun nesnelerinin state'lerinin güncellendiği fonksiyondur. Ebitengine varsayılan olarak `Update` fonksiyonunu saniyede 60 kere çağırır. Bu değer `ebiten.SetTPS` fonksiyonu çağırılarak değiştirilebilir.
- `Draw` fonkisyonu oyun nesnelerinin oyun ekranına çizdirildiği fonksiyondur. Bu fonksiyonun saniyede kaç kere çağırıldığı kullanılan ekranın `HZ` değerine göre belirlenir. Örneğin 60 Hz'lik görüntülemede saniyede 60 kere çağırılır.

```golang
// ebiten.Game interface'ini implemente eder.
type Game struct{}

// İçerisinde Game state'ini güncelleyeceğimiz fonksiyon.
// Update fonksiyonu varsayılan olarak saniyede 60 defa çağırılır.
func (g *Game) Update() error {
    return nil
}

// Draw fonksiyonu içerisinde oyun nesneleri oyun ekranı üzerinde çizdirilir.
// 60 Hz'lik görüntülemede saniyede 60 kere çağırılır
func (g *Game) Draw(screen *ebiten.Image) {

}

// outsideWidth ve outsideHeight window'un boyutlarıdır.
// Sabit bir boyut dönebiliriz ya da outsideWidth ve outsideHeight üzerinden hesaplama yapıp değer dönebiliriz.
// screenWidth, screenHeight ile outsideWidth, outsideHeight farklı olsa bile görüntü window'a sığacak şekilde otomatik olarak ölçeklenir.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 320, 240
}

func main() {
    game := &Game{}
    // Ekran boyutu ayarlanır
    ebiten.SetWindowSize(640, 480)
    ebiten.SetWindowTitle("Your game's title")
    // Game loop'u başlatmak için ebiten.RunGame fonksiyonu çağırılır.
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
```

### Grafik

#### Kutu Çizdirme

```golang
func (g *Game) Draw(screen *ebiten.Image) { 
    // ekranın tamamını beyaz renkle doldurur. Normalde ekran siyahtır.
    // ekran üzerinde gösterilecek tüm çizdirme işlemleri screen imajı üzerine yapılır.
	screen.Fill(color.White)

    // kutu boyutu
	size := 40

    // Kutuların başlangıç x ve y koordinatı
	x := 40.
	y := 10.

    // 2 kolon 5 satır kutu çizdiriyoruz
	for row := 0; row < 2; row++ {
		for col := 0; col < 5; col++ {
            // kolonların koordinatlarını satır ve kolon numaralarına göre belirliyoruz. kutular arasında 5 birim boşluk bırakıyoruz.
			rX := x + float64(col*(size+5))
			rY := y + float64(row*(size+5))

            // ebiten.NewImage bizim için boş bir imaj oluşturur
            // bu boş imajın rengini ve kordinatlarını ayarlayıp screen imajı üzerine çizdireğiz
            // dış ve iç olmak üzere 2 dikdörtgen çizdireceğiz. dış dikdörtgenin içi gri renkte, iç dikdörtgenin içi de beyaz renkte olacak. iç dikdörtgenin boyutu, dış dikdörtgenin boyutundan 4 birim az olacak. İç dikdörtgeni, dış dikdörtgenin tam ortasına yerleştireceğiz. böylece ortası beyaz gri çerçeveli bir dikdörtgen elde etmiş olacağız
			outerRect := ebiten.NewImage(size, size)
            // dış dikdörtgenin içini gri renkle dolduruyoruz
			outerRect.Fill(color.RGBA{R: 185, G: 185, B: 185, A: 255})

            // ebiten.DrawImageOptions imajı ekrana çizerken, çizdirilecek imaj üzerinde pozisyon (Translate), öteleme (Scale), döndürme (Rotate) gibi ayarlamalar yapmamızı sağlar.
			outerOpt := &ebiten.DrawImageOptions{}
            // dış dikdörtgeni rX, rY koordinatlarına öteliyoruz
			outerOpt.GeoM.Translate(rX, rY)

            // dış dikdörtgeni ekrana çizdiriyoruz
			screen.DrawImage(outerRect, outerOpt)

            // iç dikdörtgenin boyutunu dış dikdörtgenin boyutundan 4 birim küçük olarak yaratıyoruz. dış dikdörtgenin ortasına yerliştirdiğimizde kenarlarda 2 birim gri çerçeve oluşmuş olacak
			innerRect := ebiten.NewImage(size-4, size-4)
            // iç dikdörtgenin içini beyaz renkle dolduruyoruz
			innerRect.Fill(color.White)

            innerOpt := &ebiten.DrawImageOptions{}
            // iç dikdörtgeni, dış dikdörtgenin ortasına öteliyoruz
			innerOpt.GeoM.Translate(rX+2, rY+2)

            // iç dikdörtgeni ekrana çizdiriyoruz
			screen.DrawImage(innerRect, innerOpt)
		}
	}
}
```

### Animasyon

#### Tween

#### Pop

#### Shake

#### Flip

#### Wave-Like Bounce
