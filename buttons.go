package bring

const (
	MouseLeft = 1 << iota
	MouseMiddle
	MouseRight
	MouseUp
	MouseDown
)

type KeyCode []int

var (
	KeyAgain                = KeyCode{0xFF66}
	KeyAllCandidates        = KeyCode{0xFF3D}
	KeyAlphanumeric         = KeyCode{0xFF30}
	KeyLeftAlt              = KeyCode{0xFFE9}
	KeyRightAlt             = KeyCode{0xFFE9, 0xFE03}
	KeyAttn                 = KeyCode{0xFD0E}
	KeyAltGraph             = KeyCode{0xFE03}
	KeyArrowDown            = KeyCode{0xFF54}
	KeyArrowLeft            = KeyCode{0xFF51}
	KeyArrowRight           = KeyCode{0xFF53}
	KeyArrowUp              = KeyCode{0xFF52}
	KeyBackspace            = KeyCode{0xFF08}
	KeyCapsLock             = KeyCode{0xFFE5}
	KeyCancel               = KeyCode{0xFF69}
	KeyClear                = KeyCode{0xFF0B}
	KeyConvert              = KeyCode{0xFF21}
	KeyCopy                 = KeyCode{0xFD15}
	KeyCrsel                = KeyCode{0xFD1C}
	KeyCrSel                = KeyCode{0xFD1C}
	KeyCodeInput            = KeyCode{0xFF37}
	KeyCompose              = KeyCode{0xFF20}
	KeyLeftControl          = KeyCode{0xFFE3}
	KeyRightControl         = KeyCode{0xFFE3, 0xFFE4}
	KeyContextMenu          = KeyCode{0xFF67}
	KeyDelete               = KeyCode{0xFFFF}
	KeyDown                 = KeyCode{0xFF54}
	KeyEnd                  = KeyCode{0xFF57}
	KeyEnter                = KeyCode{0xFF0D}
	KeyEraseEof             = KeyCode{0xFD06}
	KeyEscape               = KeyCode{0xFF1B}
	KeyExecute              = KeyCode{0xFF62}
	KeyExsel                = KeyCode{0xFD1D}
	KeyExSel                = KeyCode{0xFD1D}
	KeyF1                   = KeyCode{0xFFBE}
	KeyF2                   = KeyCode{0xFFBF}
	KeyF3                   = KeyCode{0xFFC0}
	KeyF4                   = KeyCode{0xFFC1}
	KeyF5                   = KeyCode{0xFFC2}
	KeyF6                   = KeyCode{0xFFC3}
	KeyF7                   = KeyCode{0xFFC4}
	KeyF8                   = KeyCode{0xFFC5}
	KeyF9                   = KeyCode{0xFFC6}
	KeyF10                  = KeyCode{0xFFC7}
	KeyF11                  = KeyCode{0xFFC8}
	KeyF12                  = KeyCode{0xFFC9}
	KeyF13                  = KeyCode{0xFFCA}
	KeyF14                  = KeyCode{0xFFCB}
	KeyF15                  = KeyCode{0xFFCC}
	KeyF16                  = KeyCode{0xFFCD}
	KeyF17                  = KeyCode{0xFFCE}
	KeyF18                  = KeyCode{0xFFCF}
	KeyF19                  = KeyCode{0xFFD0}
	KeyF20                  = KeyCode{0xFFD1}
	KeyF21                  = KeyCode{0xFFD2}
	KeyF22                  = KeyCode{0xFFD3}
	KeyF23                  = KeyCode{0xFFD4}
	KeyF24                  = KeyCode{0xFFD5}
	KeyFind                 = KeyCode{0xFF68}
	KeyGroupFirst           = KeyCode{0xFE0C}
	KeyGroupLast            = KeyCode{0xFE0E}
	KeyGroupNext            = KeyCode{0xFE08}
	KeyGroupPrevious        = KeyCode{0xFE0A}
	KeyFullWidth            = KeyCode(nil)
	KeyHalfWidth            = KeyCode(nil)
	KeyHangulMode           = KeyCode{0xFF31}
	KeyHankaku              = KeyCode{0xFF29}
	KeyHanjaMode            = KeyCode{0xFF34}
	KeyHelp                 = KeyCode{0xFF6A}
	KeyHiragana             = KeyCode{0xFF25}
	KeyHiraganaKatakana     = KeyCode{0xFF27}
	KeyHome                 = KeyCode{0xFF50}
	KeyHyper                = KeyCode{0xFFED, 0xFFED, 0xFFEE}
	KeyInsert               = KeyCode{0xFF63}
	KeyJapaneseHiragana     = KeyCode{0xFF25}
	KeyJapaneseKatakana     = KeyCode{0xFF26}
	KeyJapaneseRomaji       = KeyCode{0xFF24}
	KeyJunjaMode            = KeyCode{0xFF38}
	KeyKanaMode             = KeyCode{0xFF2D}
	KeyKanjiMode            = KeyCode{0xFF21}
	KeyKatakana             = KeyCode{0xFF26}
	KeyLeft                 = KeyCode{0xFF51}
	KeyMeta                 = KeyCode{0xFFE7, 0xFFE7, 0xFFE8}
	KeyModeChange           = KeyCode{0xFF7E}
	KeyNumLock              = KeyCode{0xFF7F}
	KeyPageDown             = KeyCode{0xFF56}
	KeyPageUp               = KeyCode{0xFF55}
	KeyPause                = KeyCode{0xFF13}
	KeyPlay                 = KeyCode{0xFD16}
	KeyPreviousCandidate    = KeyCode{0xFF3E}
	KeyPrintScreen          = KeyCode{0xFF61}
	KeyRedo                 = KeyCode{0xFF66}
	KeyRight                = KeyCode{0xFF53}
	KeyRomanCharacters      = KeyCode(nil)
	KeyScroll               = KeyCode{0xFF14}
	KeySelect               = KeyCode{0xFF60}
	KeySeparator            = KeyCode{0xFFAC}
	KeyLeftShift            = KeyCode{0xFFE1}
	KeyRightShift           = KeyCode{0xFFE1, 0xFFE2}
	KeySingleCandidate      = KeyCode{0xFF3C}
	KeySuper                = KeyCode{0xFFEB, 0xFFEB, 0xFFEC}
	KeyTab                  = KeyCode{0xFF09}
	KeyUIKeyInputDownArrow  = KeyCode{0xFF54}
	KeyUIKeyInputEscape     = KeyCode{0xFF1B}
	KeyUIKeyInputLeftArrow  = KeyCode{0xFF51}
	KeyUIKeyInputRightArrow = KeyCode{0xFF53}
	KeyUIKeyInputUpArrow    = KeyCode{0xFF52}
	KeyUp                   = KeyCode{0xFF52}
	KeyUndo                 = KeyCode{0xFF65}
	KeyWin                  = KeyCode{0xFFEB}
	KeyZenkaku              = KeyCode{0xFF28}
	KeyZenkakuHankaku       = KeyCode{0xFF2}
)
