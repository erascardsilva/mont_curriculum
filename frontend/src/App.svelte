<script>
  import { onMount } from 'svelte';
  import { 
    GetProfile, SaveProfile, SelectPhoto, 
    GetEducation, SaveEducation, DeleteEducation,
    GetExperiences, SaveExperience, DeleteExperience,
    GetProjects, SaveProject, DeleteProject,
    ExportPDF, ImportResume, ExportDOCX, GetSettings, SaveSettings, ExportPDFFromData 
  } from '../wailsjs/go/main/App';

  import { BrowserOpenURL } from '../wailsjs/runtime/runtime';

  // Sistema de Tradução (i18n)
  let lang = 'pt';
  let template = 'modern';
  let settings = { language: 'pt', template: 'modern', show_photo: true, labels: {} };
  let showAbout = false; // Novo estado para o modal

  function handleDonate() {
    BrowserOpenURL('https://www.paypal.com/ncp/payment/8V6WQCGN6HDCQ');
  }

  const translations = {
    en: {
      appName: 'Monte Curriculum',
      dashSubtitle: 'Create or import your resume in minutes with premium design.',
      createNew: 'Create New',
      createNewDesc: 'Start a resume from scratch step by step.',
      import: 'Import',
      importDesc: 'Import data from an existing PDF or Word file.',
      step1: 'Personal Information',
      step2: 'Professional Objective',
      step3: 'Academic Background',
      step4: 'Work Experience',
      step5: 'Projects',
      step6: 'Templates',
      step7: 'Review and Export',
      name: 'First Name',
      lastName: 'Last Name',
      email: 'Email',
      phone: 'Phone',
      location: 'Location',
      address: 'Address',
      objective: 'Objective',
      education: 'Education',
      experience: 'Experience',
      projects: 'Projects',
      add: '+ Add New',
      cancel: 'Cancel',
      confirm: 'Confirm',
      prev: 'Previous',
      next: 'Next',
      finish: 'Finish',
      exportPdf: 'Export PDF',
      exportDocx: 'Export DOCX',
      importNotice: 'Check and adjust the imported data below.',
      templatesTitle: 'Resume Layout',
      ready: 'All Ready!',
      generating: 'Generating...',
      success: 'Exported successfully!',
      deleteConfirm: 'Are you sure you want to delete this item?',
      aboutTitle: 'About the System',
      aboutDesc: 'Monte Curriculum is an engineering solution for automation and refinement of technical resumes. The system uses advanced text processing algorithms and high-fidelity layouts to transform raw data into impeccable professional documents.'
    },
    pt: {
      appName: 'Monte Curriculum',
      dashSubtitle: 'Crie ou importe seu currículo em minutos com design premium.',
      createNew: 'Criar Novo',
      createNewDesc: 'Comece um currículo do zero passo a passo.',
      import: 'Importar',
      importDesc: 'Importe dados de um PDF ou Word existente.',
      step1: 'Informações Pessoais',
      step2: 'Objetivo Profissional',
      step3: 'Formação Acadêmica',
      step4: 'Experiência Profissional',
      step5: 'Projetos',
      step6: 'Modelos',
      step7: 'Revisão e Exportação',
      name: 'Nome',
      lastName: 'Sobrenome',
      email: 'E-mail',
      phone: 'Telefone',
      location: 'Localização',
      address: 'Endereço',
      objective: 'Objetivo',
      education: 'Formação',
      experience: 'Experiência',
      projects: 'Projetos',
      add: '+ Adicionar',
      cancel: 'Cancelar',
      confirm: 'Confirmar',
      prev: 'Anterior',
      next: 'Próximo',
      finish: 'Finalizar',
      exportPdf: 'Exportar PDF',
      exportDocx: 'Exportar DOCX',
      importNotice: 'Verifique e ajuste os dados importados abaixo.',
      templatesTitle: 'Layout do Currículo',
      ready: 'Tudo Pronto!',
      generating: 'Gerando...',
      success: 'Exportado com sucesso!',
      deleteConfirm: 'Tem certeza que deseja excluir este item?',
      aboutTitle: 'Sobre o Sistema',
      aboutDesc: 'O Monte Curriculum é uma solução de engenharia para automação e refinamento de currículos técnicos. O sistema utiliza algoritmos avançados de processamento de texto e layouts de alta fidelidade para transformar dados brutos em documentos profissionais impecáveis.'
    }
  };

  let t;
  $: t = (key) => {
    return translations[lang] && translations[lang][key] ? translations[lang][key] : key;
  };

  let view = 'dashboard'; 
  let currentStep = 1;
  const totalSteps = 7;

  let profile = { id: 1, first_name: '', last_name: '', email: '', phone: '', address: '', age: 0, photo: '', objective: '', linkedin: '', github: '', website: '' };
  let education = [];
  let experiences = [];
  let projects = [];

  let showAddEdu = false;
  let showAddExp = false;
  let showAddProj = false;

  let newEdu = { id: 0, institution: '', course: '', start_date: '', end_date: '', description: '' };
  let newExp = { id: 0, company: '', position: '', start_date: '', end_date: '', description: '' };
  let newProj = { id: 0, name: '', description: '', url: '' };

  let exportStatus = '';
  let isImportedData = false;

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    try {
      const s = await GetSettings();
      if (s) {
        settings = s;
        lang = s.language;
        template = s.template;
      }
      profile = await GetProfile();
      education = await GetEducation();
      experiences = await GetExperiences();
      projects = await GetProjects();
    } catch (e) {
      console.error("Erro ao carregar dados:", e);
    }
  }

  async function generatePDF() {
    try {
      await ExportPDFFromData(profile, education, experiences, projects, template);
    } catch (e) {
      console.error("Erro ao exportar PDF:", e);
    }
  }

  function applyFormat(field, type, id = null) {
      let target;
      if (field === 'objective') target = profile;
      else if (field === 'experience') target = experiences.find(e => e.id === id);
      else if (field === 'education') target = education.find(e => e.id === id);
      else if (field === 'project') target = projects.find(e => e.id === id);
      else if (field === 'experience_new') target = newExp;
      else if (field === 'education_new') target = newEdu;
      else if (field === 'project_new') target = newProj;

      if (!target) return;

      if (type === 'bold') {
          target.description = target.description ? target.description + ' **texto**' : '**texto**';
          if (field === 'objective') target.objective = target.objective + ' **texto**';
      } else if (type === 'bullet') {
          const prefix = target.description && !target.description.endsWith('\n') ? '\n- ' : '- ';
          if (field === 'objective') {
              profile.objective = profile.objective ? profile.objective + (profile.objective.endsWith('\n') ? '- ' : '\n- ') : '- ';
          } else {
              target.description = target.description ? target.description + prefix : '- ';
          }
      }
      
      profile = profile;
      experiences = experiences;
      education = education;
      projects = projects;
      newEdu = newEdu;
      newExp = newExp;
      newProj = newProj;
  }

  async function updateSettings() {
    settings.language = lang;
    settings.template = template;
    await SaveSettings(settings);
  }

  $: if (lang || template) {
    updateSettings();
  }

  async function startNew() {
    profile = { 
      id: 1, first_name: '', last_name: '', email: '', phone: '', address: '', 
      age: 0, photo: '', objective: '', linkedin: '', github: '', website: '' 
    };
    education = [];
    experiences = [];
    projects = [];
    isImportedData = false;
    currentStep = 1;
    view = 'wizard';
  }

  async function handleImport() {
    try {
      const imported = await ImportResume();
      if (imported) {
        await loadData();
        isImportedData = true;
        currentStep = 1;
        view = 'wizard';
      }
    } catch (e) {
      alert("Error importing: " + e);
    }
  }

  async function handlePhotoSelect() {
    const base64 = await SelectPhoto();
    if (base64) profile.photo = base64;
  }

  function sanitizeJS(text) {
    if (!text) return "";
    let sanitized = text.replace(/\r/g, "");
    sanitized = sanitized.replace(/[ \t]{2,}/g, " ");
    sanitized = sanitized.replace(/([,.!?;:])([^\s\d])/g, "$1 $2");
    sanitized = sanitized.replace(/(\.|\?|!)\s+([a-z])/g, (match, p1, p2) => p1 + " " + p2.toUpperCase());
    sanitized = sanitized.replace(/^([a-z])/, (match, p1) => p1.toUpperCase());
    return sanitized.trim();
  }

  async function handleSaveEducation() {
    if (newEdu.institution && newEdu.course) {
      await SaveEducation(newEdu);
      await loadData();
      resetEdu();
    }
  }
  function editEdu(edu) {
    newEdu = { ...edu };
    showAddEdu = true;
  }
  async function removeEdu(id) {
    if (confirm(t('deleteConfirm'))) {
      await DeleteEducation(id);
      await loadData();
    }
  }
  function resetEdu() {
    newEdu = { id: 0, institution: '', course: '', start_date: '', end_date: '', description: '' };
    showAddEdu = false;
  }

  async function handleSaveExperience() {
    if (newExp.company && newExp.position) {
      await SaveExperience(newExp);
      await loadData();
      resetExp();
    }
  }
  function editExp(exp) {
    newExp = { ...exp };
    showAddExp = true;
  }
  async function removeExp(id) {
    if (confirm(t('deleteConfirm'))) {
      await DeleteExperience(id);
      await loadData();
    }
  }
  function resetExp() {
    newExp = { id: 0, company: '', position: '', start_date: '', end_date: '', description: '' };
    showAddExp = false;
  }

  async function handleSaveProject() {
    if (newProj.name) {
      await SaveProject(newProj);
      await loadData();
      resetProj();
    }
  }
  function editProj(proj) {
    newProj = { ...proj };
    showAddProj = true;
  }
  async function removeProj(id) {
    if (confirm(t('deleteConfirm'))) {
      await DeleteProject(id);
      await loadData();
    }
  }
  function resetProj() {
    newProj = { id: 0, name: '', description: '', url: '' };
    showAddProj = false;
  }

  async function handleExportDOCX() {
    exportStatus = t('generating');
    try {
      await ExportDOCX();
      exportStatus = t('success');
      setTimeout(() => exportStatus = '', 3000);
    } catch (e) {
      alert("Error exporting DOCX: " + e);
      exportStatus = '';
    }
  }

  async function nextStep() {
    if (currentStep < totalSteps) {
      await SaveProfile(profile);
      currentStep++;
    }
  }

  function prevStep() {
    if (currentStep > 1) {
        currentStep--;
    } else {
        view = 'dashboard';
    }
  }

  const steps = [
    { id: 1, key: 'step1' },
    { id: 2, key: 'step2' },
    { id: 3, key: 'step3' },
    { id: 4, key: 'step4' },
    { id: 5, key: 'step5' },
    { id: 6, key: 'step6' },
    { id: 7, key: 'step7' }
  ];

  const templates = [
    { id: 'modern', name: 'Moderno Blue', color: '#38bdf8', previewStyle: 'grid-template-columns: 30% 70%' },
    { id: 'professional', name: 'Professional Slate', color: '#64748b', previewStyle: 'flex-direction: column' },
    { id: 'vibrant', name: 'Vibrant Orange', color: '#fb923c', previewStyle: 'grid-template-columns: 70% 30%' }
  ];
</script>

<main class="main-container">
  <div class="lang-selector">
    <button class="btn-about" on:click={() => showAbout = true}>{lang === 'pt' ? 'Sobre' : 'About'}</button>
    <button class:active={lang === 'pt'} on:click={() => lang = 'pt'} aria-label="Português">🇧🇷</button>
    <button class:active={lang === 'en'} on:click={() => lang = 'en'} aria-label="English">🇺🇸</button>
    
    <button class="btn-donate-top ripple" on:click={handleDonate} title={lang === 'pt' ? 'Apoie o projeto' : 'Support the project'}>
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M14.06 3.713c.12-1.071-.093-1.832-.702-2.526C12.628.356 11.312 0 9.626 0H4.734a.7.7 0 0 0-.691.59L2.005 13.509a.42.42 0 0 0 .415.486h2.756l-.202 1.28a.628.628 0 0 0 .62.726H8.14c.429 0 .793-.31.862-.731l.025-.13.48-3.043.03-.164.001-.007a.35.35 0 0 1 .348-.297h.38c1.266 0 2.425-.256 3.345-.91q.57-.403.993-1.005a4.94 4.94 0 0 0 .88-2.195c.242-1.246.13-2.356-.57-3.154a2.7 2.7 0 0 0-.76-.59l-.094-.061ZM6.543 8.82a.7.7 0 0 1 .321-.079H8.3c2.82 0 5.027-1.144 5.672-4.456l.003-.016q.326.186.548.438c.546.623.679 1.535.45 2.71-.272 1.397-.866 2.307-1.663 2.874-.802.57-1.842.815-3.043.815h-.38a.87.87 0 0 0-.863.734l-.03.164-.48 3.043-.024.13-.001.004a.35.35 0 0 1-.348.296H5.595a.106.106 0 0 1-.105-.123l.208-1.32z"/>
      </svg>
      <span>{lang === 'pt' ? 'Apoiar' : 'Support'}</span>
    </button>
  </div>

  {#if view === 'dashboard'}
    <div class="dashboard-view fade-in">
      <div class="dashboard-header">
        <h1>{t('appName')}</h1>
        <p>{t('dashSubtitle')}</p>
      </div>

      <div class="dashboard-grid">
        <button class="dash-card ripple" on:click={startNew} aria-label={t('createNew')}>
          <div class="icon">✨</div>
          <h3>{t('createNew')}</h3>
          <p>{t('createNewDesc')}</p>
        </button>

        <button class="dash-card ripple" on:click={handleImport} aria-label={t('import')}>
          <div class="icon">📂</div>
          <h3>{t('import')}</h3>
          <p>{t('importDesc')}</p>
        </button>
      </div>
    </div>

  {:else}
    <div class="wizard-view fade-in">
      <div class="wizard-header">
        <div class="header-top">
            <button class="btn-back-dash" on:click={() => view = 'dashboard'}>✕</button>
            <h1>{t(steps[currentStep-1].key)}</h1>
        </div>
        <div class="step-indicator">
          {#each steps as step}
            <div class="step-dot" class:active={currentStep >= step.id}></div>
          {/each}
        </div>
      </div>

      <div class="glass-panel">
        {#if isImportedData}
          <div class="import-notice fade-in">
              <span class="pulse">ℹ</span> {t('importNotice')}
          </div>
        {/if}

        {#if Number(currentStep) === 1}
          <div class="profile-layout">
            <div class="photo-section">
              <button type="button" class="photo-placeholder" on:click={handlePhotoSelect} aria-label="Upload Photo">
                {#if profile.photo}
                  <img src={profile.photo} alt="Profile" />
                {:else}
                  <div class="upload-icon">+</div>
                  <span>{lang === 'pt' ? 'Foto' : 'Photo'}</span>
                {/if}
              </button>
              <button class="btn btn-secondary btn-sm" on:click={handlePhotoSelect}>{lang === 'pt' ? 'Alterar Foto' : 'Change Photo'}</button>
            </div>

            <div class="form-grid">
              <div class="form-group" class:highlight={isImportedData && profile.first_name}>
                <label for="first_name">{t('name')}</label>
                <input id="first_name" type="text" bind:value={profile.first_name} />
              </div>
              <div class="form-group" class:highlight={isImportedData && profile.last_name}>
                <label for="last_name">{t('lastName')}</label>
                <input id="last_name" type="text" bind:value={profile.last_name} />
              </div>
              <div class="form-group" class:highlight={isImportedData && profile.email}>
                <label for="email">{t('email')}</label>
                <input id="email" type="email" bind:value={profile.email} />
              </div>
              <div class="form-group" class:highlight={isImportedData && profile.phone}>
                <label for="phone">{t('phone')}</label>
                <input id="phone" type="text" bind:value={profile.phone} />
              </div>
              <div class="form-group">
                <label for="address">{t('address')}</label>
                <input id="address" type="text" bind:value={profile.address} />
              </div>
              <div class="form-group">
                <label for="linkedin">LinkedIn</label>
                <input id="linkedin" type="url" bind:value={profile.linkedin} placeholder="linkedin.com/in/..." />
              </div>
              <div class="form-group">
                <label for="github">GitHub</label>
                <input id="github" type="url" bind:value={profile.github} placeholder="github.com/..." />
              </div>
              <div class="form-group full-width">
                  <label for="website">{lang === 'pt' ? 'Site / Portfólio' : 'Portfolio / Website'}</label>
                  <input id="website" type="url" bind:value={profile.website} placeholder="https://..." />
              </div>
            </div>
          </div>

        {:else if Number(currentStep) === 2}
          <div class="form-group" class:highlight={isImportedData && profile.objective}>
            <div class="label-edit">
              <label for="objective">{t('objective')}</label>
              <div class="text-tools">
                <button type="button" on:click={() => applyFormat('objective', 'bold')}><b>B</b></button>
                <button type="button" on:click={() => applyFormat('objective', 'bullet')}>• Tópico</button>
              </div>
            </div>
            <textarea 
              id="objective" 
              bind:value={profile.objective} 
              rows="10"
              on:blur={() => profile.objective = sanitizeJS(profile.objective)}
              placeholder={lang === 'pt' ? 'Resuma sua trajetória e objetivos...' : 'Summarize your trajectory and goals...'}
            ></textarea>
          </div>

        {:else if Number(currentStep) === 3}
          <div class="list-section">
            <div class="section-header">
              <h2>{t('education')}</h2>
              <button class="btn btn-primary btn-sm" on:click={() => showAddEdu = !showAddEdu}>
                {showAddEdu ? t('cancel') : t('add')}
              </button>
            </div>
            {#if showAddEdu}
              <div class="add-container glass-panel">
                <div class="form-grid">
                  <div class="form-group">
                    <label for="inst">{lang === 'pt' ? 'Instituição' : 'Institution'}</label>
                    <input id="inst" type="text" bind:value={newEdu.institution} on:blur={() => newEdu.institution = sanitizeJS(newEdu.institution)} />
                  </div>
                  <div class="form-group">
                    <label for="course">{lang === 'pt' ? 'Curso' : 'Course'}</label>
                    <input id="course" type="text" bind:value={newEdu.course} on:blur={() => newEdu.course = sanitizeJS(newEdu.course)} />
                  </div>
                  <div class="form-group">
                    <label for="edu_start">{lang === 'pt' ? 'Data Inicial' : 'Start Date'}</label>
                    <input id="edu_start" type="text" bind:value={newEdu.start_date} on:blur={() => newEdu.start_date = sanitizeJS(newEdu.start_date)} placeholder="Ex: 2018" />
                  </div>
                  <div class="form-group">
                    <label for="edu_end">{lang === 'pt' ? 'Data Final' : 'End Date'}</label>
                    <input id="edu_end" type="text" bind:value={newEdu.end_date} on:blur={() => newEdu.end_date = sanitizeJS(newEdu.end_date)} placeholder="Ex: 2022" />
                  </div>
                </div>
                <button class="btn btn-primary" on:click={handleSaveEducation}>{t('confirm')}</button>
              </div>
            {/if}

            <div class="items-list">
              {#each (education || []) as edu}
                <div class="item-card fade-in">
                  <div class="item-content">
                    <strong>{edu.course}</strong> — {edu.institution}
                    <small style="display:block; color:var(--primary-color)">
                      {edu.start_date} {edu.end_date ? '— ' + edu.end_date : ''}
                    </small>
                  </div>
                  <div class="item-actions">
                    <button class="btn-icon" on:click={() => editEdu(edu)} title="Editar">✏️</button>
                    <button class="btn-icon" on:click={() => removeEdu(edu.id)} title="Excluir">🗑️</button>
                  </div>
                </div>
              {/each}
            </div>
          </div>

        {:else if Number(currentStep) === 4}
           <div class="list-section">
            <div class="section-header">
              <h2>{t('experience')}</h2>
              <button class="btn btn-primary btn-sm" on:click={() => showAddExp = !showAddExp}>
                {showAddExp ? t('cancel') : t('add')}
              </button>
            </div>
            {#if showAddExp}
               <div class="add-container glass-panel">
                <div class="form-grid">
                  <div class="form-group">
                    <label for="comp">{lang === 'pt' ? 'Empresa' : 'Company'}</label>
                    <input id="comp" type="text" bind:value={newExp.company} on:blur={() => newExp.company = sanitizeJS(newExp.company)} />
                  </div>
                  <div class="form-group">
                    <label for="pos">{lang === 'pt' ? 'Cargo' : 'Title'}</label>
                    <input id="pos" type="text" bind:value={newExp.position} on:blur={() => newExp.position = sanitizeJS(newExp.position)} />
                  </div>
                  <div class="form-group">
                    <label for="exp_start">{lang === 'pt' ? 'Data Inicial' : 'Start Date'}</label>
                    <input id="exp_start" type="text" bind:value={newExp.start_date} placeholder="Ex: Jan/2020" />
                  </div>
                  <div class="form-group">
                    <label for="exp_end">{lang === 'pt' ? 'Data Final' : 'End Date'}</label>
                    <input id="exp_end" type="text" bind:value={newExp.end_date} disabled={newExp.end_date === 'Atual' || newExp.end_date === 'Present'} placeholder="Ex: Out/2022" />
                  </div>
                  <div class="form-group full-width">
                    <div class="checkbox-group" on:click={() => {
                      if (newExp.end_date === 'Atual' || newExp.end_date === 'Present') {
                        newExp.end_date = '';
                      } else {
                        newExp.end_date = lang === 'pt' ? 'Atual' : 'Present';
                      }
                    }}>
                      <input type="checkbox" id="current_work" checked={newExp.end_date === 'Atual' || newExp.end_date === 'Present'} />
                      <label for="current_work">{lang === 'pt' ? 'Trabalho Atual' : 'Current Work'}</label>
                    </div>
                  </div>
                  <div class="form-group full-width">
                    <div class="label-edit">
                      <label for="exp_desc">{lang === 'pt' ? 'Descrição das Atividades' : 'Job Description'}</label>
                      <div class="text-tools">
                        <button type="button" on:click={() => applyFormat('experience_new', 'bold')}><b>B</b></button>
                        <button type="button" on:click={() => applyFormat('experience_new', 'bullet')}>• Tópico</button>
                      </div>
                    </div>
                    <textarea id="exp_desc" rows="5" bind:value={newExp.description} on:blur={() => newExp.description = sanitizeJS(newExp.description)}></textarea>
                  </div>
                </div>
                <button class="btn btn-primary" on:click={handleSaveExperience}>{t('confirm')}</button>
              </div>
            {/if}

            <div class="items-list">
              {#each (experiences || []) as exp}
                <div class="item-card fade-in">
                  <div class="item-content">
                    <strong>{exp.position}</strong> @ {exp.company}
                    <small style="display:block; color:var(--primary-color)">
                      {exp.start_date} {exp.end_date ? '— ' + exp.end_date : ''}
                    </small>
                  </div>
                  <div class="item-actions">
                    <button class="btn-icon" on:click={() => editExp(exp)} title="Editar">✏️</button>
                    <button class="btn-icon" on:click={() => removeExp(exp.id)} title="Excluir">🗑️</button>
                  </div>
                </div>
              {/each}
            </div>
          </div>

        {:else if Number(currentStep) === 5}
           <div class="list-section">
            <div class="section-header">
              <h2>{t('projects')}</h2>
              <button class="btn btn-primary btn-sm" on:click={() => showAddProj = !showAddProj}>
                {showAddProj ? t('cancel') : t('add')}
              </button>
            </div>
              {#if showAddProj}
                <div class="add-container glass-panel">
                 <div class="form-grid">
                   <div class="form-group full-width">
                       <label for="pname">{lang === 'pt' ? 'Nome do Projeto' : 'Project Name'}</label>
                       <input id="pname" type="text" bind:value={newProj.name} on:blur={() => newProj.name = sanitizeJS(newProj.name)} />
                   </div>
                   <div class="form-group full-width">
                       <label for="purl">{lang === 'pt' ? 'Link do Projeto (URL)' : 'Project Link (URL)'}</label>
                       <input id="purl" type="url" bind:value={newProj.url} placeholder="https://..." />
                   </div>
                   <div class="form-group full-width">
                     <div class="label-edit">
                       <label for="pdesc">{lang === 'pt' ? 'Descrição' : 'Description'}</label>
                       <div class="text-tools">
                         <button type="button" on:click={() => applyFormat('project_new', 'bold')}><b>B</b></button>
                         <button type="button" on:click={() => applyFormat('project_new', 'bullet')}>• Tópico</button>
                       </div>
                     </div>
                     <textarea 
                       id="pdesc" 
                       rows="5" 
                       bind:value={newProj.description}
                       on:blur={() => newProj.description = sanitizeJS(newProj.description)}
                     ></textarea>
                   </div>
                 </div>
                 <button class="btn btn-primary" on:click={handleSaveProject}>{t('confirm')}</button>
               </div>
             {/if}
            <div class="items-list">
              {#each (projects || []) as proj}
                <div class="item-card fade-in">
                  <div class="item-content">
                    <strong>{proj.name}</strong>
                    {#if proj.url}<small style="display:block; color:var(--primary-color)">{proj.url}</small>{/if}
                    {#if proj.description}<p style="font-size:0.8rem; margin:0.5rem 0">{proj.description}</p>{/if}
                  </div>
                   <div class="item-actions">
                    <button class="btn-icon" on:click={() => editProj(proj)} title="Editar">✏️</button>
                    <button class="btn-icon" on:click={() => removeProj(proj.id)} title="Excluir">🗑️</button>
                  </div>
                </div>
              {/each}
            </div>
          </div>

        {:else if Number(currentStep) === 6}
          <div class="template-selector">
            <h2>{t('templatesTitle')}</h2>
            <div class="template-grid">
              {#each templates as t}
                <button 
                  type="button"
                  class="template-card ripple" 
                  class:selected={template === t.id}
                  on:click={() => template = t.id}
                >
                  <div class="template-mockup" style="{t.previewStyle}; border: 2px solid {template === t.id ? t.color : 'transparent'}">
                    <div class="mock-header" style="height: 10%; background: {t.color}; opacity: 0.2"></div>
                    <div class="mock-body" style="display: flex; height: 90%; gap: 2px; padding: 2px;">
                        <div class="mock-side" style="width: 30%; background: #f8fafc; opacity: 0.5; border-radius: 1px;"></div>
                        <div class="mock-main" style="width: 70%; display: flex; flex-direction: column; gap: 2px;">
                            <div style="height: 4px; background: #e2e8f0; width: 80%"></div>
                            <div style="height: 2px; background: #f1f5f9; width: 60%"></div>
                            <div style="height: 2px; background: #f1f5f9; width: 90%"></div>
                        </div>
                    </div>
                  </div>
                  <span style="color: {template === t.id ? t.color : '#fff'}">{t.name}</span>
                </button>
              {/each}
            </div>
          </div>

        {:else}
          <div class="review-section">
            <h2>{t('ready')}</h2>
            <div class="export-display">
              <button class="btn btn-primary btn-large" on:click={generatePDF}>
                <svg viewBox="0 0 24 24" width="20" height="20"><path fill="currentColor" d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"/></svg>
                {lang === 'pt' ? 'Exportar PDF Premium' : 'Export Premium PDF'}
              </button>
              <button class="btn btn-secondary btn-lg" on:click={handleExportDOCX} disabled={exportStatus !== ''}>
                {exportStatus || t('exportDocx')}
              </button>
            </div>
            {#if exportStatus}
              <p class="status-msg">{exportStatus}</p>
            {/if}
          </div>
        {/if}
      </div>

      <div class="navigation-btns">
        <button class="btn btn-secondary" on:click={prevStep}>
          {t('prev')}
        </button>
        <button class="btn btn-primary" on:click={nextStep} style:display={Number(currentStep) === totalSteps ? 'none' : 'block'}>
          {t('next')}
        </button>
      </div>
    </div>
  {/if}

  {#if showAbout}
    <div class="modal-overlay fade-in" on:click|self={() => showAbout = false}>
      <div class="modal-content glass-panel slide-up">
        <button class="modal-close" on:click={() => showAbout = false}>✕</button>
        <h2>{t('aboutTitle')}</h2>
        <p>{t('aboutDesc')}</p>
        
        <div class="donation-modal-section">
          <p class="donate-text">{lang === 'pt' ? 'Este app é completo e gratuito, ajude o desenvolvedor com uma doação:' : 'This app is full and free, please support the developer with a donation:'}</p>
          <button class="btn-donate ripple" on:click={handleDonate}>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M14.06 3.713c.12-1.071-.093-1.832-.702-2.526C12.628.356 11.312 0 9.626 0H4.734a.7.7 0 0 0-.691.59L2.005 13.509a.42.42 0 0 0 .415.486h2.756l-.202 1.28a.628.628 0 0 0 .62.726H8.14c.429 0 .793-.31.862-.731l.025-.13.48-3.043.03-.164.001-.007a.35.35 0 0 1 .348-.297h.38c1.266 0 2.425-.256 3.345-.91q.57-.403.993-1.005a4.94 4.94 0 0 0 .88-2.195c.242-1.246.13-2.356-.57-3.154a2.7 2.7 0 0 0-.76-.59l-.094-.061ZM6.543 8.82a.7.7 0 0 1 .321-.079H8.3c2.82 0 5.027-1.144 5.672-4.456l.003-.016q.326.186.548.438c.546.623.679 1.535.45 2.71-.272 1.397-.866 2.307-1.663 2.874-.802.57-1.842.815-3.043.815h-.38a.87.87 0 0 0-.863.734l-.03.164-.48 3.043-.024.13-.001.004a.35.35 0 0 1-.348.296H5.595a.106.106 0 0 1-.105-.123l.208-1.32z"/>
            </svg>
            {lang === 'pt' ? 'Apoie o projeto' : 'Support this project'}
          </button>
        </div>

        <div class="modal-signature">
           <div class="sig-name">Erasmo Cardoso</div>
           <div class="sig-role">Software Engineer | Electronics Techniciant</div>
        </div>
      </div>
    </div>
  {/if}
</main>

<style>
  .lang-selector {
    position: absolute;
    top: 2.5rem;
    right: 3rem;
    display: flex;
    align-items: center;
    gap: 1.5rem;
    z-index: 100;
  }

  .lang-selector button {
    background: rgba(255,255,255,0.05);
    border: 1px solid var(--border-color);
    font-size: 1.5rem;
    padding: 0.5rem;
    border-radius: 0.5rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .lang-selector button.active {
    background: var(--primary-color);
    border-color: var(--primary-color);
    transform: scale(1.1);
  }

  .btn-about {
    font-size: 0.75rem !important;
    text-transform: uppercase;
    font-weight: 700;
    letter-spacing: 1.5px;
    padding: 0.6rem 1.2rem !important;
    color: var(--text-color);
    border: 1px solid var(--border-color) !important;
    background: rgba(255,255,255,0.05) !important;
    border-radius: 2rem;
    cursor: pointer;
    transition: all 0.2s;
    height: fit-content;
    line-height: 1;
  }

  .btn-about:hover {
    color: var(--primary-color);
    background: rgba(56, 189, 248, 0.1) !important;
    border-color: var(--primary-color) !important;
    transform: translateY(-2px);
  }

  .btn-donate {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    font-size: 0.9rem !important;
    text-transform: uppercase;
    font-weight: 700;
    letter-spacing: 1px;
    padding: 0.8rem 1.5rem !important;
    color: #fff;
    border: 1px solid #1e3a8a !important;
    background: linear-gradient(135deg, #0070ba 0%, #003087 100%) !important;
    border-radius: 2.5rem;
    cursor: pointer;
    transition: all 0.2s;
    height: fit-content;
    line-height: 1;
    box-shadow: 0 4px 12px rgba(0, 112, 186, 0.3);
    margin: 0 auto;
  }

  .donation-modal-section {
    margin-bottom: 3rem;
    padding: 1.5rem;
    background: rgba(255,255,255,0.03);
    border-radius: 1rem;
    border: 1px solid rgba(255,255,255,0.05);
  }

  .donate-text {
    font-size: 0.9rem !important;
    margin-bottom: 1.2rem !important;
    color: var(--text-muted) !important;
  }

  .modal-signature {
    margin-top: 1rem;
    padding-top: 1.5rem;
    border-top: 1px solid rgba(255,255,255,0.05);
  }

  .sig-name {
    font-size: 1.2rem;
    color: var(--primary-color);
    font-weight: 700;
    margin-bottom: 0.2rem;
  }

  .sig-role {
    font-size: 0.85rem;
    color: var(--text-muted);
    font-weight: 500;
  }

  /* Modal Styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(15, 23, 42, 0.85);
    backdrop-filter: blur(12px);
    z-index: 2000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  .modal-content {
    max-width: 550px;
    width: 100%;
    padding: 3.5rem;
    position: relative;
    border: 1px solid rgba(255,255,255,0.1);
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.7);
    text-align: center;
    background: rgba(30, 41, 59, 0.7);
  }

  .modal-close {
    position: absolute;
    top: 1.2rem;
    right: 1.2rem;
    background: rgba(255,255,255,0.05);
    border: 1px solid rgba(255,255,255,0.1);
    color: var(--text-muted);
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s;
  }

  .modal-close:hover {
    color: #fff;
    background: rgba(255,255,255,0.1);
    transform: rotate(90deg);
  }

  /* Text Tools Styling */
  .text-tools {
    display: flex;
    gap: 0.5rem;
  }

  .text-tools button {
    background: rgba(255,255,255,0.05);
    border: 1px solid var(--border-color);
    color: var(--text-color);
    padding: 0.2rem 0.6rem;
    border-radius: 4px;
    font-size: 0.75rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .text-tools button:hover {
    background: rgba(56, 189, 248, 0.1);
    border-color: var(--primary-color);
    color: var(--primary-color);
  }

  .slide-up {
    animation: slideUp 0.5s cubic-bezier(0.16, 1, 0.3, 1);
  }

  @keyframes slideUp {
    from { transform: translateY(40px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }

  /* Dashboard Styles */
  .dashboard-view {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 60vh;
    text-align: center;
    max-width: 900px;
    margin: 2rem auto;
  }

  .wizard-view {
    display: flex;
    flex-direction: column;
    padding-bottom: 2rem;
  }

  .dashboard-header h1 {
    font-size: 4rem;
    margin-bottom: 1rem;
    background: linear-gradient(135deg, #fff 0%, #38bdf8 100%);
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .dashboard-header p {
    font-size: 1.25rem;
    color: var(--text-secondary);
    margin-bottom: 3rem;
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    width: 100%;
  }

  .dash-card {
    background: var(--surface-color);
    border: 1px solid var(--border-color);
    border-radius: 2rem;
    padding: 3rem 2rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    color: #fff;
  }

  .dash-card:hover {
    transform: translateY(-10px);
    border-color: var(--primary-color);
    background: rgba(56, 189, 248, 0.05);
  }

  .dash-card h3 { 
    color: #ffffff !important; 
    font-size: 1.5rem;
    margin: 1rem 0 0.5rem 0;
  }

  .dash-card p {
    color: #cbd5e1 !important;
  }

  .dash-card .icon { font-size: 3rem; }

  /* Wizard Styles */
  .main-container {
    padding: 2rem;
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    position: relative;
    overflow-y: visible;
  }

  .wizard-header { margin-bottom: 2rem; }
  .header-top { display: flex; align-items: center; gap: 1.5rem; margin-bottom: 1rem; }
  .btn-back-dash { background: none; border: none; color: var(--text-secondary); font-size: 1.5rem; cursor: pointer; }

  .profile-layout { display: flex; gap: 2.5rem; align-items: flex-start; }
  .photo-section { display: flex; flex-direction: column; align-items: center; gap: 1.5rem; flex-shrink: 0; }
  .photo-placeholder {
    width: 160px; height: 160px; border-radius: 1.5rem; background: var(--surface-color); border: 2px dashed var(--border-color);
    cursor: pointer; overflow: hidden; position: relative; display: flex; justify-content: center; align-items: center;
  }
  .photo-placeholder img { width: 100%; height: 100%; object-fit: cover; }

  .form-grid { flex: 1; display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  .full-width { grid-column: span 2; }
  .form-group label { display: block; margin-bottom: 0.5rem; color: var(--text-secondary); font-size: 0.9rem; }
  
  .glass-panel { padding: 2.5rem; margin-bottom: 2rem; }

  .navigation-btns { 
    display: flex; 
    justify-content: space-between; 
    margin-top: 2rem; 
    padding: 1.5rem 0;
    border-top: 1px solid var(--border-color);
  }

  .item-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: rgba(255,255,255,0.03);
    padding: 1rem;
    border-radius: 1rem;
    border: 1px solid var(--border-color);
    margin-bottom: 0.5rem;
  }

  .item-actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn-icon {
    background: none;
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0.25rem;
    transition: transform 0.2s;
  }

  .btn-icon:hover {
    transform: scale(1.2);
  }

  /* Template Selection Styles */
  .template-selector {
    text-align: center;
    padding: 1rem 0;
  }

  .template-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 2.5rem;
    padding: 2rem 0;
    max-width: 1000px;
    margin: 0 auto;
  }

  .template-card {
    background: var(--surface-color);
    border: 1px solid var(--border-color);
    border-radius: 1.5rem;
    padding: 1.5rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    position: relative;
    overflow: hidden;
  }

  .template-card:hover {
    transform: translateY(-8px);
    border-color: var(--primary-color);
    background: rgba(56, 189, 248, 0.05);
  }

  .template-card.selected {
    border-color: var(--primary-color);
    background: rgba(56, 189, 248, 0.1);
    box-shadow: 0 0 30px rgba(56, 189, 248, 0.2);
  }

  .template-mockup {
    width: 100%;
    aspect-ratio: 1 / 1.4;
    background: #fff;
    border-radius: 0.5rem;
    box-shadow: 0 10px 20px rgba(0,0,0,0.2);
    transition: transform 0.3s;
  }

  .template-card:hover .template-mockup {
    transform: scale(1.05);
  }

  .template-card span {
    font-weight: 600;
    font-size: 1.1rem;
    letter-spacing: 0.5px;
  }

  @media (max-width: 768px) {
    .dashboard-grid { grid-template-columns: 1fr; }
    .profile-layout { flex-direction: column; align-items: center; }
    .template-grid { grid-template-columns: 1fr; }
  }

  .btn-donate-top {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    background: linear-gradient(135deg, #0070ba 0%, #003087 100%);
    color: white;
    border: 1px solid rgba(255,255,255,0.1);
    padding: 0.5rem 1rem;
    border-radius: 2rem;
    font-size: 0.8rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 1px;
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(0, 48, 135, 0.3);
    transition: all 0.3s ease;
  }

  .btn-donate-top:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 15px rgba(0, 48, 135, 0.5);
    filter: brightness(1.1);
  }

  .btn-donate-top svg {
    filter: drop-shadow(0 1px 2px rgba(0,0,0,0.2));
  }

  .checkbox-group {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
    user-select: none;
    padding: 0.5rem 0;
  }
  .checkbox-group input {
    width: 1.2rem;
    height: 1.2rem;
    cursor: pointer;
    accent-color: var(--primary-color);
  }
  .checkbox-group label {
    margin: 0 !important;
    cursor: pointer;
    color: var(--text-primary) !important;
  }
</style>
