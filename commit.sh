#!/bin/bash

# Verifica se o primeiro argumento é "tag"
if [ "$1" = "tag" ]; then
  CREATE_TAG="s"
else
  CREATE_TAG="n"
fi

# Verifica se go.work existe antes de renomear
if [ -f go.work ]; then
  # Passo 1: Renomear go.work para go.work.backup
  mv go.work go.work.backup
  WORK_RENAMED="yes"
else
  WORK_RENAMED="no"
fi

# Passo 2: Limpar o cache do Go
go clean -modcache

# Atualizar dependências específicas
echo "Atualizando dependências..."

# Passo 3: Executar um go mod tidy
go mod tidy

# Passo 4: Executar o commit com a mensagem passada como segundo argumento, se disponível
if [ -n "$2" ]; then
  COMMIT_MESSAGE=$2
else
  echo "Digite a mensagem de commit:"
  read COMMIT_MESSAGE
fi

git add .
git commit -m "$COMMIT_MESSAGE"

# Passo 5: Sincronizar os fontes
git push

# Restaurar o go.work apenas se foi renomeado anteriormente neste script
if [ "$WORK_RENAMED" = "yes" ]; then
  mv go.work.backup go.work 2>/dev/null
fi

# Criar e fazer push de uma tag, se requisitado
if [ "$CREATE_TAG" = "s" ]; then
  # Buscar a última tag
  LAST_TAG=$(git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null)

  if [ -z "$LAST_TAG" ]; then
    # Se não houver nenhuma tag anterior, usar v0.0.1
    NEW_TAG="v0.0.1"
  else
    # Incrementar a tag. Assume o formato v0.0.x.
    NEW_TAG=$(echo $LAST_TAG | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g')
  fi

  # Criar a nova tag
  git tag $NEW_TAG

  # Fazer push da tag
  git push origin $NEW_TAG
fi
