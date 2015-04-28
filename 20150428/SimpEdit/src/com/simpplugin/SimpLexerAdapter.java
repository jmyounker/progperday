package com.simpplugin;

/**
 * Created by jeff on 4/28/15.
 */

import com.intellij.lexer.FlexAdapter;
import java.io.Reader;

public class SimpLexerAdapter extends FlexAdapter {
    public SimpLexerAdapter() {
        super(new SimpLexer((Reader) null));
    }
}
