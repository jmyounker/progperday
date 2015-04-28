package com.simpplugin.psi;

/**
 * Created by jeff on 4/28/15.
 */
import com.intellij.extapi.psi.PsiFileBase;
import com.intellij.openapi.fileTypes.FileType;
import com.intellij.psi.FileViewProvider;
import com.simpplugin.SimpFileType;
import com.simpplugin.SimpLanguage;
import org.jetbrains.annotations.NotNull;

import javax.swing.*;

public class SimpFile extends PsiFileBase {
    public SimpFile(@NotNull FileViewProvider viewProvider) {
        super(viewProvider, SimpLanguage.INSTANCE);
    }

    @NotNull
    @Override
    public FileType getFileType() {
        return SimpFileType.INSTANCE;
    }

    @Override
    public String toString() {
        return "Simp File";
    }

    @Override
    public Icon getIcon(int flags) {
        return super.getIcon(flags);
    }
}
