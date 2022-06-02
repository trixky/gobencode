package parser

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestParseBytes(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// empty
		{input: "0:", expected: ""},
		// plain
		{input: "1:o", expected: "o"},
		{input: "4:chat", expected: "chat"},
		// long
		{input: "16:0123456789abcdef", expected: "0123456789abcdef"},
	}

	for _, test := range tests {
		output, err := parseBytes(bufio.NewReader(strings.NewReader(test.input[1:])), byte(test.input[0]))

		if err != nil {
			t.Errorf("failed to parse input [%s]: %v\n", test.input, err)
		}

		if reflect.TypeOf(output) != reflect.TypeOf(test.expected) {
			t.Errorf("bad output format: output [%T] != [%T] (expected)\n", output, test.expected)
			continue
		}

		if output != test.expected {
			t.Errorf("output [%s] != [%s] expected\n", output, test.expected)
		}
	}
}

func TestParseInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// null
		{input: "i0e", expected: 0}, // TEST LE -0 ?
		// positive
		{input: "i1e", expected: 1},
		{input: "i9e", expected: 9},
		{input: "i123e", expected: 123},
		// negative
		{input: "i-1e", expected: -1},
		{input: "i-9e", expected: -9},
		{input: "i-123e", expected: -123},
	}

	for _, test := range tests {
		output, err := parseInteger(bufio.NewReader(strings.NewReader(test.input[1:])))

		if err != nil {
			t.Errorf("failed to parse input [%s]: %v\n", test.input, err)
			continue
		}

		if reflect.TypeOf(output) != reflect.TypeOf(test.expected) {
			t.Errorf("bad output format: output [%T] != [%T] (expected)\n", output, test.expected)
			continue
		}

		if output != test.expected {
			t.Errorf("output %d != %d (expected)\n", output, test.expected)
		}
	}
}

func TestParseList(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// empty
		{input: "le", expected: []interface{}{}},
		// only integers
		{input: "li0ee", expected: []interface{}{0}},
		{input: "li1ee", expected: []interface{}{1}},
		{input: "li-1ee", expected: []interface{}{-1}},
		{input: "li1ei-2ei3ee", expected: []interface{}{1, -2, 3}},
		// only strings
		{input: "l0:e", expected: []interface{}{""}},
		{input: "l3:ouie", expected: []interface{}{"oui"}},
		{input: "l3:oui3:non2:oke", expected: []interface{}{"oui", "non", "ok"}},
		//  integers and strings
		{input: "li1e3:ouii2e3:noni3e2:oke", expected: []interface{}{1, "oui", 2, "non", 3, "ok"}},
		//  list in list
		{input: "llelelee", expected: []interface{}{[]interface{}{}, []interface{}{}, []interface{}{}}},
		{input: "lli1eeli2eeli3eee", expected: []interface{}{[]interface{}{1}, []interface{}{2}, []interface{}{3}}},
		{input: "ll3:ouiel3:nonel2:okee", expected: []interface{}{[]interface{}{"oui"}, []interface{}{"non"}, []interface{}{"ok"}}},
		//  dictionary in list
		{input: "ldededee", expected: []interface{}{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}}},
		{input: "ld3:oui3:noned2:oki2eed2:koi0eee", expected: []interface{}{map[string]interface{}{"oui": "non"}, map[string]interface{}{"ok": 2}, map[string]interface{}{"ko": 0}}},
	}

	for _, test := range tests {
		output, err := parseList(bufio.NewReader(strings.NewReader(test.input[1:])))

		if err != nil {
			t.Errorf("failed to parse input [%s]: %v\n", test.input, err)
			continue
		}

		if reflect.TypeOf(output) != reflect.TypeOf(test.expected) {
			t.Errorf("bad output format: output [%T] != [%T] (expected)\n", output, test.expected)
			continue
		}

		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("output %v != %v (expected)\n", output, test.expected)
		}
	}
}

func TestParseDictionaryionary(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// empty
		{input: "de", expected: map[string]interface{}{}},
		// only integers
		{input: "d0:i0ee", expected: map[string]interface{}{"": 0}},
		{input: "d4:zeroi0ee", expected: map[string]interface{}{"zero": 0}},
		{input: "d8:moins uni-1ee", expected: map[string]interface{}{"moins un": -1}},
		{input: "d5:troisi3e8:moins uni-1ee", expected: map[string]interface{}{"trois": 3, "moins un": -1}},
		// only strings
		{input: "d0:0:e", expected: map[string]interface{}{"": ""}},
		{input: "d0:3:ouie", expected: map[string]interface{}{"": "oui"}},
		{input: "d3:non3:ouie", expected: map[string]interface{}{"non": "oui"}},
		{input: "d3:non3:oui4:cinq7:derniere", expected: map[string]interface{}{"non": "oui", "cinq": "dernier"}},
		{input: "d3:non3:oui4:cinq7:dernier4:vide0:e", expected: map[string]interface{}{"non": "oui", "cinq": "dernier", "vide": ""}},
		//  integers and strings
		{input: "d0:i0e1:x1:ye", expected: map[string]interface{}{"": 0, "x": "y"}},
		{input: "d0:i0e1:x1:y8:moins uni-1ee", expected: map[string]interface{}{"": 0, "x": "y", "moins un": -1}},
		//  dictionary in dictionary
		{input: "d2:abde2:cdde2:efdee", expected: map[string]interface{}{"ab": map[string]interface{}{}, "cd": map[string]interface{}{}, "ef": map[string]interface{}{}}},
		{input: "d2:abd3:ouii3ee2:cdd3:noni4ee2:efd3:popi5eee", expected: map[string]interface{}{"ab": map[string]interface{}{"oui": 3}, "cd": map[string]interface{}{"non": 4}, "ef": map[string]interface{}{"pop": 5}}},
		//  list in dictionary
		{input: "d2:able2:cdle2:eflee", expected: map[string]interface{}{"ab": []interface{}{}, "cd": []interface{}{}, "ef": []interface{}{}}},
		{input: "d2:abl3:oui4:chate2:cdle2:efl5:chienee", expected: map[string]interface{}{"ab": []interface{}{"oui", "chat"}, "cd": []interface{}{}, "ef": []interface{}{"chien"}}},
	}

	for _, test := range tests {
		output, err := parseDictionary(bufio.NewReader(strings.NewReader(test.input[1:])))

		if err != nil {
			t.Errorf("failed to parse input [%s]: %v\n", test.input, err)
			continue
		}

		if reflect.TypeOf(output) != reflect.TypeOf(test.expected) {
			t.Errorf("bad output format: output [%T] != [%T] (expected)\n", output, test.expected)
			continue
		}

		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("output %v != %v (expected)\n", output, test.expected)
		}
	}
}

func TestParseElementBasicTorrent(t *testing.T) {
	tests := []string{
		`d8:announce35:https://torrent.ubuntu.com/announce13:announce-listll35:https://torrent.ubuntu.com/announceel40:https://ipv6.torrent.ubuntu.com/announceee7:comment29:Ubuntu CD releases.ubuntu.com10:created by13:mktorrent 1.113:creation datei1650550976e4:infod6:lengthi3654957056e4:name30:ubuntu-22.04-desktop-amd64.iso12:piece lengthi262144e6:pieces3:...ee`,
		`d8:announce59:dht://B98537878AE9F32F647EB4913D649C2DB1C337F1.dht/announce13:announce-listll59:dht://B98537878AE9F32F647EB4913D649C2DB1C337F1.dht/announceel44:udp://tracker.openbittorrent.com:80/announceel38:udp://tracker.publicbt.com:80/announceel34:udp://tracker.coppersurfer.tk:6969el40:udp://tracker.leechers-paradise.org:6969el27:udp://open.demonii.com:1337ee18:azureus_propertiesd7:Contentd15:Content Networki3e11:Description512:"Every conversation has the potential to open up and reveal all the layers and layers within it, all those rooms within rooms," says podcaster and musician Hrishikesh Hirway. In this profoundly moving talk, he offers a guide to deep conversations and explores what you learn when you stop to listen closely. Stay tuned to the end to hear a performance of his original song "Between There and Here (feat. Yo-Yo Ma)."<img src="http://feeds.feedburner.com/~r/TedtalksHD/~4/WYZlXqvdRd8" height="1" width="1" alt=""/>8:Durationi913000e9:Thumbnail3:...eee`,
		`d8:announce44:udp://tracker.openbittorrent.com:80/announce13:announce-listll44:udp://tracker.openbittorrent.com:80/announceel42:udp://tracker.opentrackr.org:1337/announceel38:udp://tracker.publicbt.com:80/announceee10:created by19:qBittorrent v3.3.1513:creation datei1612723546e4:infod5:filesld6:lengthi19710976e4:pathl13:resource1.bineed6:lengthi9050674e4:pathl12:resource.bineed6:lengthi1056768e4:pathl9:setup.exeeee4:name16:Minecraft 1.15.212:piece lengthi16384e6:pieces3:...ee`,
		`d8:announce35:https://torrent.ubuntu.com/announce7:comment29:Kubuntu CD cdimage.ubuntu.com10:created by13:mktorrent 1.113:creation datei1650550153e4:infod6:lengthi3674746880e4:name31:kubuntu-22.04-desktop-amd64.iso12:piece lengthi262144e6:pieces3:...ee`,
		`d8:announce59:dht://ABF6A85BAAEBEE731A79C1BAF83326F5C70F76CB.dht/announce13:announce-listll59:dht://ABF6A85BAAEBEE731A79C1BAF83326F5C70F76CB.dht/announceel44:udp://tracker.openbittorrent.com:80/announceel38:udp://tracker.publicbt.com:80/announceel34:udp://tracker.coppersurfer.tk:6969el40:udp://tracker.leechers-paradise.org:6969el27:udp://open.demonii.com:1337ee18:azureus_propertiesd7:Contentd15:Content Networki3e12:Content Type8:featured11:Description36:A Seriously Bloody Disgusting Bundle9:Thumbnail3:...eee`,
		`d8:announce59:dht://B8FAD0E6B49E7A33542A393FE75884E829C22706.dht/announce13:announce-listll59:dht://B8FAD0E6B49E7A33542A393FE75884E829C22706.dht/announceel44:udp://tracker.openbittorrent.com:80/announceel38:udp://tracker.publicbt.com:80/announceel34:udp://tracker.coppersurfer.tk:6969el40:udp://tracker.leechers-paradise.org:6969el27:udp://open.demonii.com:1337ee18:azureus_propertiesd7:Contentd15:Content Networki3e11:Description434:Journalist Lara Bitar on economic and political collapse in Lebanon; An Afghan interpreter who rescued Joe Biden in 2008 escapes the Taliban with help from private groups; Correspondent for The Nation John Nichols on Donald Trump's attempts to overturn the election and the January 6 coup at Capitol Hill that nearly succeeded. Get Democracy Now! delivered right to your inbox. Sign up for the Daily Digest: democracynow.org/subscribe8:Durationi3540000e9:Thumbnail3:...eee`,
		`d7:comment41:Arch Linux 2022.05.01 (www.archlinux.org)10:created by13:mktorrent 1.113:creation datei1651396816e4:infod6:lengthi866463744e4:name31:archlinux-2022.05.01-x86_64.iso12:piece lengthi524288e6:pieces3:...ee`,
		`d8:announce14:/announce?pid=13:announce-listll42:udp://tracker.opentrackr.org:1337/announceel33:udp://open.stealth.si:80/announceel32:udp://explodie.org:6969/announceel37:udp://exodus.desync.com:6969/announceel48:udp://tracker.internetwarriors.net:1337/announceel39:udp://ipv4.tracker.harry.lu:80/announceel30:udp://9.rarbg.to:2740/announceel31:udp://9.rarbg.com:2770/announceel31:http://1337.abcvg.info/announceel41:udp://tracker.torrent.eu.org:451/announceel37:http://tracker.bt4g.com:2095/announceel31:udp://opentor.org:2710/announceel37:udp://www.torrent.eu.org:451/announceel42:udp://retracker.lanta-net.ru:2710/announceee10:created by37:ruTorrent (PHP Class - Adrien Gibrat)13:creation datei1653416473e4:infod6:lengthi755813766e4:name53:[ Torrent911.net ] Code 252 - Signal de détresse.avi12:piece lengthi524288e6:pieces3:3eeee`,
		`d8:announce14:/announce?pid=13:announce-listll42:udp://tracker.opentrackr.org:1337/announceel33:udp://open.stealth.si:80/announceel32:udp://explodie.org:6969/announceel37:udp://exodus.desync.com:6969/announceel48:udp://tracker.internetwarriors.net:1337/announceel39:udp://ipv4.tracker.harry.lu:80/announceel30:udp://9.rarbg.to:2740/announceel31:udp://9.rarbg.com:2770/announceel31:http://1337.abcvg.info/announceel41:udp://tracker.torrent.eu.org:451/announceel37:http://tracker.bt4g.com:2095/announceel31:udp://opentor.org:2710/announceel37:udp://www.torrent.eu.org:451/announceel42:udp://retracker.lanta-net.ru:2710/announceee10:created by30:Transmission/2.94 (d8e60ee44f)13:creation datei1642172247e8:encoding5:UTF-84:infod6:lengthi728921018e4:name58:[ Torrent911.com ] Gaia.2021.FRENCH.HDRip.XviD-EXTREME.avi12:piece lengthi524288e6:pieces3:...eee`,
		`d8:announce14:/announce?pid=13:announce-listll42:udp://tracker.opentrackr.org:1337/announceel33:udp://open.stealth.si:80/announceel32:udp://explodie.org:6969/announceel37:udp://exodus.desync.com:6969/announceel48:udp://tracker.internetwarriors.net:1337/announceel39:udp://ipv4.tracker.harry.lu:80/announceel30:udp://9.rarbg.to:2740/announceel31:udp://9.rarbg.com:2770/announceel31:http://1337.abcvg.info/announceel41:udp://tracker.torrent.eu.org:451/announceel37:http://tracker.bt4g.com:2095/announceel31:udp://opentor.org:2710/announceel37:udp://www.torrent.eu.org:451/announceel42:udp://retracker.lanta-net.ru:2710/announceee10:created by30:Transmission/2.94 (d8e60ee44f)13:creation datei1653417512e8:encoding5:UTF-84:infod6:lengthi253537285e4:name95:[ Torrent911.net ] Le.Flambeau.Les.Aventuriers.de.Chupacabra.S01E03.FRENCH.WEB.x264-EXTREME.mkv12:piece lengthi131072e6:pieces3:...ee`,
		`d8:announce14:/announce?pid=13:announce-listll42:udp://tracker.opentrackr.org:1337/announceel33:udp://open.stealth.si:80/announceel32:udp://explodie.org:6969/announceel37:udp://exodus.desync.com:6969/announceel48:udp://tracker.internetwarriors.net:1337/announceel39:udp://ipv4.tracker.harry.lu:80/announceel30:udp://9.rarbg.to:2740/announceel31:udp://9.rarbg.com:2770/announceel31:http://1337.abcvg.info/announceel41:udp://tracker.torrent.eu.org:451/announceel37:http://tracker.bt4g.com:2095/announceel31:udp://opentor.org:2710/announceel37:udp://www.torrent.eu.org:451/announceel42:udp://retracker.lanta-net.ru:2710/announceee10:created by37:ruTorrent (PHP Class - Adrien Gibrat)13:creation datei1653418540e4:infod6:lengthi1526549551e4:name72:[ Torrent911.net ] Pourris.Gâtés.2021.FRENCH.720p.WEB.x264-EXTREME.mkv12:piece lengthi524288e6:pieces3:...ee`,
		`d8:announce14:/announce?pid=13:announce-listll42:udp://tracker.opentrackr.org:1337/announceel33:udp://open.stealth.si:80/announceel32:udp://explodie.org:6969/announceel37:udp://exodus.desync.com:6969/announceel48:udp://tracker.internetwarriors.net:1337/announceel39:udp://ipv4.tracker.harry.lu:80/announceel30:udp://9.rarbg.to:2740/announceel31:udp://9.rarbg.com:2770/announceel31:http://1337.abcvg.info/announceel41:udp://tracker.torrent.eu.org:451/announceel37:http://tracker.bt4g.com:2095/announceel31:udp://opentor.org:2710/announceel37:udp://www.torrent.eu.org:451/announceel42:udp://retracker.lanta-net.ru:2710/announceee10:created by30:Transmission/2.94 (d8e60ee44f)13:creation datei1642191432e8:encoding5:UTF-84:infod6:lengthi1475481942e4:name80:[ Torrent911.com ] The.Tragedy.of.Macbeth.2021.TRUEFRENCH.HDRip.XviD-EXTREME.avi12:piece lengthi1048576e6:pieces3:...ee`,
		`d8:announce14:/announce?pid=13:announce-listll42:udp://tracker.opentrackr.org:1337/announceel33:udp://open.stealth.si:80/announceel32:udp://explodie.org:6969/announceel37:udp://exodus.desync.com:6969/announceel48:udp://tracker.internetwarriors.net:1337/announceel39:udp://ipv4.tracker.harry.lu:80/announceel30:udp://9.rarbg.to:2740/announceel31:udp://9.rarbg.com:2770/announceel31:http://1337.abcvg.info/announceel41:udp://tracker.torrent.eu.org:451/announceel37:http://tracker.bt4g.com:2095/announceel31:udp://opentor.org:2710/announceel37:udp://www.torrent.eu.org:451/announceel42:udp://retracker.lanta-net.ru:2710/announceee10:created by30:Transmission/2.94 (d8e60ee44f)13:creation datei1653415595e8:encoding5:UTF-84:infod6:lengthi557985820e4:name70:[ Torrent911.net ] We.Own.This.City.S01E05.FRENCH.WEB.x264-EXTREME.mkv12:piece lengthi524288e6:pieces3:...ee`,
	}

	for index, test := range tests {
		_, err := ParseElement(bufio.NewReader(strings.NewReader(test)))

		if err != nil {
			t.Errorf("failed to parse input %d [%s...]: %v\n", index+1, test[:50], err)

			return
		}
	}
}
