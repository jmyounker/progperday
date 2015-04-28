package com.simpplugin;

import com.intellij.lang.Language;

/**
 * Created by jeff on 4/24/15.
 */
public class SimpLanguage extends Language {
    public static final SimpLanguage INSTANCE = new SimpLanguage();

    private SimpLanguage() {
        super("Simp");
    }
}

