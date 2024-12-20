package main

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"pdf_greyhat_go/api"
	"pdf_greyhat_go/download"
)

func main() {
	// Initialize session and keywords
	sessionCookie := "54e7fe8c2aa1dd504b9be39fa3466f10"
	keywords := []string{
		"app",
		"main",
		"index",
		"script",
		"core",
		"init",
		"bundle",
		"config",
		"setup",
		"service",
		"utils",
		"utility",
		"functions",
		"helper",
		"helpers",
		"base",
		"vendor",
		"custom",
		"common",
		"global",
		"lib",
		"library",
		"data",
		"api",
		"client",
		"server",
		"auth",
		"login",
		"dashboard",
		"home",
		"router",
		"routes",
		"navigation",
		"validate",
		"validation",
		"form",
		"forms",
		"ajax",
		"request",
		"handler",
		"event",
		"events",
		"dom",
		"loader",
		"lazyload",
		"analytics",
		"tracking",
		"error",
		"errors",
		"debug",
		"log",
		"logs",
		"monitor",
		"polyfill",
		"worker",
		"webpack.config",
		"rollup.config",
		"gulpfile",
		"build",
		"deploy",
		"test",
		"tests",
		"spec",
		"mock",
		"mocks",
		"report",
		"reports",
	}
	extensions := map[string][]string{
		"go":   {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
		"php":  {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
		"java": {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
		"js":   {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
		"json": {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
		"py":   {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
		"xml":  {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
		"pdf":  {".*.vic.gov.au.*", ".*010-strafrechtadvocaten.nl.*", ".*010strafrecht.com.*", ".*010strafrecht.nl.*", ".*010strafrechtadvocaten.nl.*", ".*020xiangji.com.*", ".*023printking.com.*", ".*043web.nl.*", ".*0day.rocks.*", ".*0ktbv.com.*", ".*0nyx.net.*", ".*0patch.com.*", ".*0wx.cat.*", ".*0wx.es.*", ".*0wx.eu.*", ".*0wx.net.*", ".*0wx.org.*", ".*0x1d107.xyz.*", ".*0x42.sh.*", ".*0x44sec.com.*", ".*0x6761.ga.*", ".*0xdeadbeef.network.*", ".*0xff.se.*", ".*1-1cloud.com.*", ".*1-jc-handyman.com.*", ".*10-14.nl.*", ".*1000mercis.com.*", ".*1000mercis.fr.*", ".*1000zgarb.cz.*", ".*1001-mondes.fr.*", ".*100percentpure.sk.*", ".*100tb.com.*", ".*100x100art.ru.*", ".*1014onderwijs.nl.*", ".*101caffe.cz.*", ".*101media.nl.*", ".*109fsg.com.*", ".*10jahreabfuck.de.*", ".*10k.sk.*", ".*10kb.de.*", ".*10uur.nl.*", ".*10xgenomics.com.*", ".*112maashorst.nl.*", ".*117.cz.*", ".*1177.se.*", ".*11freunde.de.*", ".*11urss.com.*", ".*1232karma.com.*", ".*12345pro.vip.*", ".*1234help.nl.*", ".*123accu.nl.*", ".*123billeje.dk.*", ".*123carhire.co.uk.*", ".*123carhire.eu.*", ".*123carhire.uk.*", ".*123carrental.com.*", ".*123carrentals.com.*", ".*123d-juupbackes.nl.*", ".*123greetings.com.*", ".*123kocarky.cz.*", ".*123kociky.sk.*", ".*123led.sk.*", ".*123licence.cz.*", ".*123recht.de.*", ".*123tcs.io.*", ".*123toner.cz.*", ".*123zing.nl.*", ".*125ofsteam.com.*", ".*12pd.com.*", ".*12pointdesign.com.*", ".*12steps.io.*", ".*1337r00t.me.*", ".*139.162.176.152.*", ".*139.162.225.134.*", ".*13th-dover.uk.*", ".*140-let.cz.*", ".*14c-on.top.*", ".*172.104.98.170.*", ".*172.105.17.61.*", ".*172.105.232.137.*", ".*172.105.73.195.*", ".*17avolemsaberlaveritat.cat.*", ".*17cpis.com.*", ".*18street.com.*", ".*1907sweets.nl.*", ".*198.58.105.17.*", ".*1982.ml.*", ".*19peaks.nl.*", ".*19thcenturyscience.org.*", ".*1Password.*", ".*1Shoppingcart.com.*", ".*1arocket.de.*", ".*1bite.nl.*", ".*1blu.de.*", ".*1cdky.com.*", ".*1cloud.at.*", ".*1don-store.com.*", ".*1e.com.*", ".*1eco.nl.*", ".*1hattem.nl.*", ".*1hk.fi.*", ".*1js.de.*", ".*1kampen.nl.*", ".*1km.co.in.*", ".*1limburg.nl.*", ".*1me.cz.*", ".*1milliontruths.com.*", ".*1mobil.sk.*", ".*1msvelkeopatovice.cz.*", ".*1obuv.cz.*", ".*1op1dieetsintmichielsgestel.nl.*", ".*1password.ca.*", ".*1password.com.*", ".*1password.eu.*", ".*1pstage.com.*", ".*1rs.nl.*", ".*1se.co.*", ".*1teddy.cz.*", ".*1und1.de.*", ".*1undbesser.de.*", ".*1xbet-aze.com.*", ".*1xbet-bo.com.*", ".*1xbet-sw.com.*", ".*1xbet.com.*", ".*1xbet.com.ve.*", ".*1xbet.gp.*", ".*1xbet82.com.*", ".*1xlite-529281.top.*", ".*1xmailjhwk.top.*", ".*1xsportcom.com.*", ".*1xsultan.com.*", ".*1zwartewaterland.nl.*", ".*1zwolle.nl.*", ".*2-way.nl.*", ".*2020fidelity.info.*", ".*2020saving.info.*", ".*2061.club.*", ".*20forma.nl.*", ".*20muleteamlaundry.com.*", ".*21cdc.co.*", ".*220triathlon.com.*", ".*233info.tk.*", ".*23andMe.*", ".*23andballing.com.*", ".*23crows.co.uk.*", ".*23g.io.*", ".*24257.eu.*", ".*24fit.cz.*", ".*24hero.nl.*", ".*24heures.ch.*", ".*24hourallservices.com.*", ".*24idcheck.com.*", ".*24liveresults.com.*", ".*24uurloodgieternederland.nl.*", ".*24uurtaxiderondevenen.nl.*", ".*25.wf.*", ".*2501.xyz.*", ".*25space.com.*", ".*25up.co.uk.*", ".*2718282.net.*", ".*27pro-portal.de.*", ".*2chan.eu.*", ".*2chan.jp.*", ".*2cleangutter.com.*", ".*2connect.nl.*", ".*2cu.nu.*", ".*2cutedecoration.be.*", ".*2degrees.nz.*", ".*2dhype.nl.*", ".*2din.cz.*", ".*2dinradia.cz.*", ".*2do.best.*", ".*2ehandsgoudensieraden.nl.*", ".*2gdnc.org.*", ".*2girlsandasign.com.*", ".*2gosoftware.eu.*", ".*2j2k.cz.*", ".*2keep.net.*", ".*2lg.notaires.fr.*", ".*2lrn4.com.*", ".*2makeitwork.nl.*", ".*2manydots.nl.*", ".*2miracles.be.*", ".*2mparts.com.*", ".*2ns.fi.*", ".*2pizza.ca.*", ".*2pleasure.cz.*", ".*2pt0.org.*", ".*2sacademy.nl.*", ".*2sigma.school.*", ".*2xleshop.at.*", ".*3-en-uno.lat.*", ".*3-in-one.lat.*", ".*3-m.cz.*", ".*333obra.com.br.*", ".*337ericksen.com.*", ".*343mfsidfhie.ga.*", ".*34victorhugo.notaires.fr.*", ".*360-cb.com.*", ".*360amigo.com.*"},
	}

	createOutputFile := func(keyword string) (string, error) {
		filename := fmt.Sprintf("results-%s.json", keyword)
		dir, err := os.Open(".")
		if err != nil {
			return "", fmt.Errorf("failed opening the directory: %w", err)
		}
		defer dir.Close()

		var acc int
		names, err := dir.Readdirnames(-1)
		if err != nil && err != io.EOF { // EOF means end of directory
			return "", fmt.Errorf("error reading directory: %w", err)
		}

		for _, name := range names {
			if name == filename || name == fmt.Sprintf("results-%s-%d.json", keyword, acc) {
				acc++
			}
		}

		if acc > 0 {
			filename = fmt.Sprintf("results-%s-%d.json", keyword, acc)
		}

		return filename, nil
	}

	// ---------------------------------------------------------------
	for _, keyword := range keywords {
		outputFile, err := createOutputFile(keyword)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			continue
		}
		fmt.Printf("Searching for files with keyword: %s\n", keyword)
		var files []api.FileInfo
		maxRetries := 3
		for retries := 0; retries < maxRetries; retries++ {
			files, err = api.QueryFiles(sessionCookie, []string{keyword}, extensions)
			if err == nil {
				break
			}
			log.Printf("Retry %d/%d for keyword '%s' failed: %v", retries+1, maxRetries, keyword, err)
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			log.Printf("All retries failed for keyword '%s'\n", keyword)
			continue
		}

		var wg sync.WaitGroup
		results := make([]map[string]interface{}, 0)
		var mutex sync.Mutex

		concurrencyLimit := 6
		semaphore := make(chan struct{}, concurrencyLimit)
		errorschan := make(chan error, len(files))

		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				mutex.Lock()
				err := saveResults(results, outputFile)
				if err != nil {
					errorschan <- fmt.Errorf("error saving periodic results for keyword '%s': %v", keyword, err)
				}
				mutex.Unlock()
			}
		}()
		done := make(chan struct{})
		defer close(done)
		for _, fileInfo := range files {
			wg.Add(1)

			fmt.Println("Processing file:", fileInfo.Filename)
			if fileInfo.Size > 50*1024*1024 {
				errorschan <- fmt.Errorf("skipping large file: %s", fileInfo.Filename)
				continue
			}

			go func(file api.FileInfo) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				result := download.ProcessFile(file, extensions)
				if result != nil {
					mutex.Lock()
					results = append(results, result)
					mutex.Unlock()
					if err != nil {
						errorschan <- fmt.Errorf("failed to open file '%s' for writing: %v", outputFile, err)
						return
					}
				} else {
					errorschan <- fmt.Errorf("processing failed for file: %s", file.URL)
				}
			}(fileInfo)
		}

		go func() {
			for err := range errorschan {
				fmt.Printf("Error: %v\n", err)
			}
			close(errorschan) // Close the error channel after all errors are collected
		}()

		go func() {
			for {
				select {
				case err, ok := <-errorschan:
					if !ok {
						break
					}
					fmt.Printf("Error: %v\n", err)
				case _, ok := <-done:
					if !ok {
						break
					}
					fmt.Println("All files processed")
					return
				}
			}
		}()
		wg.Wait()
		mutex.Lock()
		err = saveResults(results, outputFile)
		if err != nil {
			log.Printf("Error saving final results for keyword '%s': %v", keyword, err)
		}
		mutex.Unlock()
	}
}

func saveResults(results []map[string]interface{}, outputFile string) error {
	fmt.Println("Saving results...")

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputFile, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to write JSON to file '%s': %w", outputFile, err)
	}

	fmt.Printf("Results saved to %s\n", outputFile)
	return nil
}
